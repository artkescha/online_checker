package router

import (
	"github.com/artkescha/checker/online_checker/pkg/middlewares"
	"github.com/artkescha/checker/online_checker/pkg/session"
	task_handlers "github.com/artkescha/checker/online_checker/web/task/handlers"
	try_handlers "github.com/artkescha/checker/online_checker/web/try/handlers"
	"github.com/artkescha/checker/online_checker/web/user/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(userHandlers handlers.User, taskHandlers task_handlers.TaskHandler,
	tryHandler try_handlers.SolutionHandler, sessionManager session.Manager) *mux.Router {

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
	router.HandleFunc("/state", userHandlers.State).Methods("GET")

	//redirect to user TODO maybe tasks (GET)
	router.HandleFunc("/user", /*middlewares.Authorization(sessionManager,*/ userHandlers.List/*)*/).Methods("GET")

	//redirect to user TODO maybe tasks (GET)
	router.HandleFunc("/admin", /*middlewares.Authorization(sessionManager, */taskHandlers.List/*)*/).Methods("GET")

	//task new
	router.HandleFunc("/tasks/new", taskHandlers.CreateForm).Methods("GET")

	//task new
	router.HandleFunc("/tasks", taskHandlers.Create).Methods("POST")

	//read one task
	router.HandleFunc("/tasks/{id}", taskHandlers.ReadOne).Methods("GET")

	//edit
	router.HandleFunc("/tasks/edit/{id}", taskHandlers.Edit).Methods("GET")

	//updated
	router.HandleFunc("/tasks/{id}", taskHandlers.Update).Methods("POST")

	//upload tests
	router.HandleFunc("/tests/upload/{taskID}", taskHandlers.UploadTests).Methods("POST")

	//download tests
	router.HandleFunc("/tests/download/{taskID}", taskHandlers.DownloadTests).Methods("GET")

	router.HandleFunc("/tasks/{id}", taskHandlers.Delete).Methods("DELETE")

	//solution form {id - taskID}
	router.HandleFunc("/tasks/solutionForm/{taskID}", taskHandlers.SolutionForm).Methods("GET")

	//send solution
	router.HandleFunc("/try", middlewares.Authorization(sessionManager, tryHandler.SendSolution)).Methods("POST")

	//read solution
	router.HandleFunc("/try/{ID}", tryHandler.ReadOneTry).Methods("GET")

	//get tries by user
	router.HandleFunc("/tries/userID/{userID}", tryHandler.ListByUserID).Methods("GET")

	//подключаем статику к форме login-а
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/user/template/"))))

	//подключаем статику к форме login-а
	router.PathPrefix("/").Handler(http.StripPrefix("/tries/userID/{userID}/", http.FileServer(http.Dir("./web/try/template/"))))

	//подключаем общие middlewares
	router.Use(middlewares.Logger.AccessLogMiddleware)

	//отсутствие урла обрабатыватся отдельно
	//router.NotFoundHandler = http.HandlerFunc(template.ExcecLogin)
	return router
}
