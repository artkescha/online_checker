package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/session"
	"github.com/artkescha/checker/online_checker/pkg/task"
	task_repo "github.com/artkescha/checker/online_checker/pkg/task/repository"
	"github.com/artkescha/checker/online_checker/pkg/user"
	"github.com/artkescha/checker/online_checker/pkg/user/repository"
	"github.com/artkescha/checker/online_checker/web/request"
	"github.com/artkescha/checker/online_checker/web/response"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type User interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	LogOut(w http.ResponseWriter, r *http.Request)
	State(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	RegisterForm(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	Tmpl           *template.Template
	UsersRepo      repository.UserRepo
	TasksRepo      task_repo.TaskRepo
	SessionManager session.Manager
	Logger         *zap.SugaredLogger
}

// регистрация нового пользователя
func (h UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := request.Login{}
	err := request.DecodePostParams(&params, r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.UsersRepo.Insert(params.Username, user.GetMD5Password(params.Password))
	if err != nil {
		paramsErrors := response.NewResponseParamsErrors("body",
			"username", params.Username, err)
		response.WriteParamsErrors(w, http.StatusUnprocessableEntity, paramsErrors)
		return
	}
	token, err := h.SessionManager.CreateSession(*user)
	if err != nil {
		paramsErrors := response.NewResponseParamsErrors("body",
			"username", params.Username, err)
		response.WriteParamsErrors(w, http.StatusUnauthorized, paramsErrors)
		return
	}
	response.WriteResponse(w, http.StatusCreated, token, "token")
}

// аутентификация пользователя
func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	params := request.Login{}
	err := request.DecodePostParams(&params, r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.UsersRepo.GetUserByLogin(params.Username)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	if !user.ValidPassword(params.Password) {
		response.WriteError(w, http.StatusUnauthorized, errors.New("invalid password"))
		return
	}
	token, err := h.SessionManager.CreateSession(*user)
	log.Print(token)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	//TODO позже переписать!!!!!!!!!!!!!!!!!
	data := struct {
		Token   string `json:"token"`
		IsAdmin bool   `json:"isAdmin"`
	}{
		Token:   token,
		IsAdmin: user.IsAdmin(),
	}
	response_, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %s", err)
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response_)
	if err != nil {
		log.Printf("write response failed: %s", err)
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//response.WriteResponse(w, http.StatusOK, token, "token")

	//-----------------------------------------------------------

	//r.Header.Set("Authorization", "Bearer " + token)
	//http.Redirect(w, r, "/state", http.StatusSeeOther)
}

func (h UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	user, err := request.ExtractContext(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	err = h.SessionManager.DestroySession(strconv.Itoa(int(user.ID)))
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	response.WriteResponse(w, http.StatusOK, struct{}{})
}

// аутентификация пользователя
func (h UserHandler) State(w http.ResponseWriter, r *http.Request) {
	user, err := request.ExtractContext(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	if user.IsAdmin() {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h UserHandler) List(w http.ResponseWriter, r *http.Request) {
	user_, err := request.ExtractContext(r)
	if err != nil {
		h.Logger.Error("extract use request context err", err)
		http.Error(w, `extract use request context err`, http.StatusUnauthorized)
		return
	}
	//TODO limit:3 offset:0 in request
	tasks, err := h.TasksRepo.List(r.Context(), 100, 0, "created_at")
	if err != nil {
		h.Logger.Error("get tasks list err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "list.html", struct {
		Tasks  []task.Task
		UserId int64
	}{
		Tasks:  tasks,
		UserId: user_.ID,
	})
	if err != nil {
		h.Logger.Error("tasks list executeTemplate err", err)
		http.Error(w, fmt.Sprintf(`tasks list template err %s`, err), http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) RegisterForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
