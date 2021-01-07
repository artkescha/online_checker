package handlers

import (
	"encoding/json"
	"errors"
	"gitlab.com/artkescha/grader/online_checker/pkg/session"
	"gitlab.com/artkescha/grader/online_checker/pkg/user"
	"gitlab.com/artkescha/grader/online_checker/pkg/user/repository"
	"gitlab.com/artkescha/grader/online_checker/web/request"
	"gitlab.com/artkescha/grader/online_checker/web/response"
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
	User(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	RegisterForm(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	Tmpl           *template.Template
	UsersRepo      repository.UserRepo
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
	w.Write(response_)

	//response.WriteResponse(w, http.StatusOK, token, "token")

	//-----------------------------------------------------------

	//r.Header.Set("Authorization", "Bearer " + token)
	//http.Redirect(w, r, "/state", http.StatusSeeOther)
}

func (h UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	session, err := request.ExtractContext(r)

	err = h.SessionManager.DestroySession(strconv.Itoa(int(session.ID)))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
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

func (h UserHandler) User(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "user.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
