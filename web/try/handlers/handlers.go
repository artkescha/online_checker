package handlers

import (
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/tries/status"
	"github.com/artkescha/checker/online_checker/pkg/writer"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/artkescha/checker/online_checker/pkg/session"
	"github.com/artkescha/checker/online_checker/pkg/tries"
	"github.com/artkescha/checker/online_checker/pkg/tries/repository"
	"github.com/artkescha/checker/online_checker/pkg/tries/transmitter"
	"github.com/artkescha/checker/online_checker/web/request"
	"github.com/artkescha/checker/online_checker/web/response"
	"github.com/artkescha/grader_api/send_solution"
)

type Solutioner interface {
	SendSolution(w http.ResponseWriter, r *http.Request)
	ListByUserID(w http.ResponseWriter, r *http.Request)
	ReadOneTry(w http.ResponseWriter, r *http.Request)
}

type SolutionHandler struct {
	Tmpl           *template.Template
	TriesRepo      repository.TriesRepo
	SessionManager session.Manager
	Transmitter    transmitter.Transmitter
	Writer         writer.Writer
	Logger         *zap.SugaredLogger
}

func (h SolutionHandler) SendSolution(w http.ResponseWriter, r *http.Request) {
	var try try.Try
	//r.ParseForm()
	err := request.DecodePostParams(&try, r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := request.ExtractContext(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	try.Created = time.Now()
	try.Status = status.Queue
	try.Description = try.Status.String()
	//try.UserID = user.ID
	log.Printf("try: %v", try)

	tryId, err := h.Writer.Write(try)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	apiTry := send_solution.Try{
		Id:         uint64(tryId),
		UserId:     uint64(user.ID),
		Solution:   try.Solution,
		Timestamp:  time.Now().Unix(),
		TaskId:     int32(try.TaskID),
		LanguageId: int32(try.LanguageID),
	}

	err = h.Transmitter.Transmit("solution", &apiTry)
	if err != nil {
		h.Logger.Error("solution transmit failed", zap.Error(err))
		http.Error(w, `{"error": "publish to broker failed"}`, http.StatusInternalServerError)
		return
	}

	//http.Redirect(w, r, "/tries/userID/64", http.StatusSeeOther)
}

func (h SolutionHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		h.Logger.Errorf("extract request params failed: %s", err)
		http.Error(w, fmt.Sprintf(`userID not found %s`, err), http.StatusInternalServerError)
		return
	}
	triesByUser, err := h.TriesRepo.ListByUser(r.Context(), uint64(userID), 100, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Info("triesByUser: %+v", triesByUser)
	err = h.Tmpl.ExecuteTemplate(w, "list.html", struct {
		Tries []try.Try
	}{
		Tries: triesByUser,
	})
	if err != nil {
		http.Error(w, `tries template err`, http.StatusInternalServerError)
		return
	}
}

func (h SolutionHandler) ReadOneTry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		h.Logger.Errorf("extract request params failed: %s", err)
		http.Error(w, fmt.Sprintf(`try id not found in request params %s`, err), http.StatusInternalServerError)
		return
	}
	try, err := h.TriesRepo.GetByID(r.Context(), uint64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Info("try with id: %d", id, try)
	err = h.Tmpl.ExecuteTemplate(w, "try.html", try)
	if err != nil {
		http.Error(w, `try template err`, http.StatusInternalServerError)
		return
	}
}
