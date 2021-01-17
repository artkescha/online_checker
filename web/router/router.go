package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/artkescha/grader/online_checker/pkg/middlewares"
	"gitlab.com/artkescha/grader/online_checker/pkg/session"
	task_handlers "gitlab.com/artkescha/grader/online_checker/web/task/handlers"
	"gitlab.com/artkescha/grader/online_checker/web/user/handlers"
	"net/http"
)

func NewRouter(userHandlers handlers.User, taskHandlers task_handlers.TaskHandler, sessionManager session.Manager) *mux.Router {
	router := mux.NewRouter()

	//регистрация пользователя
	router.HandleFunc("/register", userHandlers.Register).Methods("POST")

	//авторизация пользователя
	router.HandleFunc("/login", userHandlers.Login).Methods("POST")

	//шаблон логина
	router.HandleFunc("/", userHandlers.Index)

	//шаблон регистрации
	router.HandleFunc("/registerForm", userHandlers.RegisterForm).Methods("GET")

	//redirect user or admin
	router.HandleFunc("/state", middlewares.Authorization(sessionManager, userHandlers.State)).Methods("GET")

	//redirect to user
	router.HandleFunc("/user", middlewares.Authorization(sessionManager, userHandlers.List)).Methods("GET")

	//redirect to user
	router.HandleFunc("/admin", taskHandlers.List).Methods("GET")

	//task new
	router.HandleFunc("/tasks/new", taskHandlers.CreateForm).Methods("GET")

	//task new
	router.HandleFunc("/tasks/new", taskHandlers.Create).Methods("POST")

	//read one task
	router.HandleFunc("/tasks/{id}", taskHandlers.ReadOne).Methods("GET")

	//edit
	router.HandleFunc("/tasks/edit/{id}", taskHandlers.Edit).Methods("GET")

	//updated
	router.HandleFunc("/tasks/{id}", taskHandlers.Update).Methods("PATH")

	router.HandleFunc("/tasks/{id}", taskHandlers.Delete).Methods("DELETE")

	//send solution {id - taskID}
	router.HandleFunc("/tasks/solution/{id}", taskHandlers.SolutionForm).Methods("GET")

	//подключаем статику к форме login-а
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/user/template/"))))

	//подключаем общие middlewares
	router.Use(middlewares.Logger.AccessLogMiddleware)

	//отсутствие урла обрабатыватся отдельно
	//router.NotFoundHandler = http.HandlerFunc(template.ExcecLogin)
	return router
}
