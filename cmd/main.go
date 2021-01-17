package main

import (
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/lib/pq"
	"gitlab.com/artkescha/grader/online_checker/pkg/middlewares"
	"gitlab.com/artkescha/grader/online_checker/pkg/session"
	task_repo "gitlab.com/artkescha/grader/online_checker/pkg/task/repository"
	user_repo "gitlab.com/artkescha/grader/online_checker/pkg/user/repository"
	"gitlab.com/artkescha/grader/online_checker/web/router"
	"gitlab.com/artkescha/grader/online_checker/web/server"
	task_handlers "gitlab.com/artkescha/grader/online_checker/web/task/handlers"
	"gitlab.com/artkescha/grader/online_checker/web/user/handlers"
	"go.uber.org/zap"
	"html/template"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	//port, exists := os.LookupEnv("PORT")
	//if !exists {
	//	zapLogger.Error("port is required")
	//	return
	//}

	port := "8080"

	zapLogger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", "0.0.0.0"),
		zap.String("port", port),
	)

	logger := zapLogger.Sugar()

	middlewares.Logger.ZapLogger = logger

	//postgreURL, exists := os.LookupEnv("DATABASE_URL")
	//if !exists {
	//	zapLogger.Error("DATABASE_URL is required")
	//	return
	//}

	postgreURL := "host=localhost port=5432 user=postgres password=postgres dbname=grader sslmode=disable"
	db, err := sql.Open("postgres", postgreURL)
	if err != nil {
		zapLogger.Error("connection to postgres failed: %s", zap.Error(err))
	}
	defer db.Close()

	mc := memcache.New("127.0.0.1:11211")

	manager := session.NewManager(mc)

	userHandlers := handlers.UserHandler{
		Tmpl:      template.Must(template.ParseGlob("./web/user/template/*")),
		UsersRepo: user_repo.NewUsersRepo(db),
		//TODO дубль подумать использовать ли интерфейс!!!!!!!!!
		TasksRepo: task_repo.NewTasksRepo(db),
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

	router := router.NewRouter(userHandlers, taskHandlers, manager)

	server.Start("0.0.0.0:"+port, router)
}
