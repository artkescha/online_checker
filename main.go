package main

import (
	"database/sql"
	"go.uber.org/zap"
	"html/template"
	"math"
	"time"

	"github.com/artkescha/checker/online_checker/config"
	"github.com/artkescha/checker/online_checker/pkg/middlewares"
	"github.com/artkescha/checker/online_checker/pkg/session"
	task_repo "github.com/artkescha/checker/online_checker/pkg/task/repository"
	try_repo "github.com/artkescha/checker/online_checker/pkg/tries/repository"
	user_repo "github.com/artkescha/checker/online_checker/pkg/user/repository"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"

	"github.com/artkescha/checker/online_checker/pkg/tries/transmitter"
	"github.com/artkescha/checker/online_checker/web/router"
	"github.com/artkescha/checker/online_checker/web/server"
	task_handlers "github.com/artkescha/checker/online_checker/web/task/handlers"
	try_handlers "github.com/artkescha/checker/online_checker/web/try/handlers"
	"github.com/artkescha/checker/online_checker/web/user/handlers"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return
	}
	logger := zapLogger.Sugar()

	//TODO replace config_path !!!
	config_path := "./etc/config.yaml"

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

	logger.Info("starting server host: %s, port: %s", config_.WebUrl.Hostname(), config_.WebUrl.Port())

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

	mc := memcache.New(config_.MemoryStorageUrl)

	manager := session.NewManager(mc)

	userHandlers := handlers.UserHandler{
		Tmpl:      template.Must(template.ParseGlob("./web/user/template/*")),
		UsersRepo: user_repo.NewUsersRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		TasksRepo:      task_repo.NewTasksRepo(db),
		SessionManager: manager,
		Logger:         logger,
	}

	taskHandlers := task_handlers.TaskHandler{
		Tmpl:      template.Must(template.ParseGlob("./web/task/template/*")),
		TasksRepo: task_repo.NewTasksRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		SessionManager: manager,
		Logger:         logger,
	}

	nc, err := nats.Connect("nats://172.16.238.13:4222",
		nats.ReconnectWait(1*time.Minute),
		nats.MaxReconnects(int(math.MaxUint32)),
		nats.ReconnectHandler(func(nc *nats.Conn) {

			logger.Info("got reconnected to host %s", nc.ConnectedUrl())
		}))
	defer func() {
		nc.Close()
	}()
	if err != nil {
		logger.Error("nats connection %s", err)
		return
	}

	solutionHandler := try_handlers.SolutionHandler{
		Tmpl:      template.Must(template.ParseGlob("./web/try/templates/*")),
		TriesRepo: try_repo.NewTriesRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		SessionManager: manager,
		Transmitter:    transmitter.New(nc),
		Logger:         logger,
	}

	router := router.NewRouter(userHandlers, taskHandlers, solutionHandler, manager)

	server.Start(config_.WebUrl.Host, router)
}
