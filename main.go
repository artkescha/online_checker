package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"html/template"
	"math"

	"github.com/artkescha/checker/online_checker/config"
	"github.com/artkescha/checker/online_checker/pkg/middlewares"
	"github.com/artkescha/checker/online_checker/pkg/session"
	task_repo "github.com/artkescha/checker/online_checker/pkg/task/repository"
	try_repo "github.com/artkescha/checker/online_checker/pkg/tries/repository"
	"github.com/artkescha/checker/online_checker/pkg/tries/transmitter"
	user_repo "github.com/artkescha/checker/online_checker/pkg/user/repository"
	"github.com/artkescha/checker/online_checker/pkg/writer"
	"github.com/artkescha/checker/online_checker/web/router"
	"github.com/artkescha/checker/online_checker/web/server"
	task_handlers "github.com/artkescha/checker/online_checker/web/task/handlers"
	try_handlers "github.com/artkescha/checker/online_checker/web/try/handlers"
	"github.com/artkescha/checker/online_checker/web/user/handlers"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return
	}

	logger := zapLogger.Sugar()

	//TODO replace config_path !!!
	config_path := "./config.yaml"

	config_, err := config.ConfigFromFile(config_path)
	if err != nil {
		logger.Error("parse config file failed: %s", err)
		return
	}

	//port, exists := os.LookupEnv("PORT")
	//if !exists {
	//	zapLogger.Error("port is required")
	//	return
	//}

	//port := "8080"

	logger.Infof("starting server host: %s, port: %s", config_.WebUrl.Hostname(), config_.WebUrl.Port())

	middlewares.Logger.ZapLogger = logger

	//postgreURL, exists := os.LookupEnv("DATABASE_URL")
	//if !exists {
	//	zapLogger.Error("DATABASE_URL is required")
	//	return
	//}

	//postgreURL := "host=localhost port=5432 user=postgres password=postgres dbname=grader sslmode=disable"
	//postgreURL := "host=172.16.238.10 port=5432 user=root password=root dbname=grader sslmode=disable"
	db, err := sql.Open(config_.DBDriver, config_.DBConnection)
	if err != nil {
		logger.Error("connection to postgres failed: %s", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("db close failed: %s", err)
			return
		}
	}()

	mc := memcache.New(config_.MemcachedUrl)

	manager := session.NewManager(mc)

	userHandlers := handlers.UserHandler{
		Tmpl:      template.Must(template.ParseGlob("../web/user/template/*")),
		UsersRepo: user_repo.NewUsersRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		TasksRepo:      task_repo.NewTasksRepo(db),
		SessionManager: manager,
		Logger:         logger,
	}

	taskHandlers := task_handlers.TaskHandler{
		Tmpl:      template.Must(template.ParseGlob("../web/task/template/*")),
		TasksRepo: task_repo.NewTasksRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		SessionManager: manager,
		Config:         config_,
		Logger:         logger,
	}

	//nc, err := nats.Connect("nats://172.16.238.13:4222",
	//	nats.ReconnectWait(1*time.Minute),
	//	nats.MaxReconnects(int(math.MaxUint32)),
	//	nats.ReconnectHandler(func(nc *nats.Conn) {
	//		logger.Info("got reconnected to host %s", nc.ConnectedUrl())
	//	}))
	//defer func() {
	//	nc.Close()
	//}()
	//if err != nil {
	//	logger.Error("nats connection %s", err)
	//	return
	//}

	// Send PINGs every 10 seconds, and fail after 5 PINGs without any response.
	brokerConnect, err := stan.Connect("test-cluster", "test-client_1",
		stan.NatsURL(config_.BrokerUrl))
	stan.Pings(60, math.MaxUint32)
	stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		logger.Errorf("Connection lost, reason: %v", reason)
	})

	defer func() {
		if err := brokerConnect.Close(); err != nil {
			logger.Errorf("connect to nats-server with url %s failed: %s", config_.BrokerUrl, err)
		}
	}()
	if err != nil {
		logger.Errorf("connect to nats-server with url %s failed: %s", config_.BrokerUrl, err)
		return
	}

	solutionHandler := try_handlers.SolutionHandler{
		Tmpl:      template.Must(template.ParseGlob("../web/try/templates/*")),
		TriesRepo: try_repo.NewTriesRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		SessionManager: manager,
		Transmitter:    transmitter.New(brokerConnect, logger),
		Writer:         writer.NewDBWriter(db),
		Logger:         logger,
	}

	router := router.NewRouter(userHandlers, taskHandlers, solutionHandler, manager)

	server.Start(config_.WebUrl.Host, router)
}
