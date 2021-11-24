package main

import (
	"database/sql"
	"github.com/artkescha/checker/online_checker/config"
	"github.com/artkescha/checker/online_checker/pkg/middlewares"
	"github.com/artkescha/checker/online_checker/pkg/session"
	task_repo "github.com/artkescha/checker/online_checker/pkg/task/repository"
	try_repo "github.com/artkescha/checker/online_checker/pkg/tries/repository"
	user_repo "github.com/artkescha/checker/online_checker/pkg/user/repository"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"math"
	"time"

	"github.com/artkescha/checker/online_checker/pkg/tries/transmitter"
	"github.com/artkescha/checker/online_checker/web/router"
	"github.com/artkescha/checker/online_checker/web/server"
	task_handlers "github.com/artkescha/checker/online_checker/web/task/handlers"
	try_handler "github.com/artkescha/checker/online_checker/web/try/handlers"
	"github.com/artkescha/checker/online_checker/web/user/handlers"
	"go.uber.org/zap"
	"html/template"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	config, err := config.ConfigFromFile("./config.yaml")
	if err!=nil {
		zapLogger.Error("parse config file failed: %s", zap.Error(err))
		return
	}

	//port, exists := os.LookupEnv("PORT")
	//if !exists {
	//	zapLogger.Error("port is required")
	//	return
	//}

	//port := "8080"

	zapLogger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", config.WebUrl.Hostname()),
		zap.String("port", config.WebUrl.Port()),
	)

	logger := zapLogger.Sugar()

	middlewares.Logger.ZapLogger = logger

	//postgreURL, exists := os.LookupEnv("DATABASE_URL")
	//if !exists {
	//	zapLogger.Error("DATABASE_URL is required")
	//	return
	//}

	//postgreURL := "host=localhost port=5432 user=postgres password=postgres dbname=grader sslmode=disable"
	//postgreURL := "host=172.16.238.10 port=5432 user=root password=root dbname=grader sslmode=disable"
	db, err := sql.Open(config.DBDriver, config.DBConnection)
	if err != nil {
		zapLogger.Error("connection to postgres failed: %s", zap.Error(err))
		return
	}
	defer db.Close()

	//mc := memcache.New("172.16.238.12:11211")
	mc := memcache.New(config.MemoryStorageUrl)

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

			zapLogger.Info("got reconnected to ",
				zap.String("host:", nc.ConnectedUrl()),
			)
		}))

	if err != nil {
		zapLogger.Error("nats connection", zap.Error(err))
		return
	}

	solutionHandler := try_handler.SolutionHandler{
		Tmpl:      template.Must(template.ParseGlob("./web/try/templates/*")),
		TriesRepo: try_repo.NewTriesRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		SessionManager: manager,
		Transmitter:    transmitter.New(nc),
		Logger:         logger,
	}

	router := router.NewRouter(userHandlers, taskHandlers, solutionHandler, manager)

	server.Start(config.WebUrl.Host, router)
}
