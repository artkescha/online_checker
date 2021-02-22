package handlers

import (
	"github.com/artkescha/grader/online_checker/pkg/session"
	"github.com/artkescha/grader/online_checker/pkg/tries"
	"github.com/artkescha/grader/online_checker/pkg/tries/repository"
	"github.com/artkescha/grader/online_checker/pkg/tries/transmitter"
	"github.com/artkescha/grader/online_checker/web/response"
	"github.com/gorilla/schema"
	"time"

	"go.uber.org/zap"
	"net/http"
)

type Solutioner interface {
	SendSolution(w http.ResponseWriter, r *http.Request)
}

type SolutionHandler struct {
	TasksRepo      repository.TriesRepo
	SessionManager session.Manager
	Transmitter    transmitter.Transmitter
	Logger         *zap.SugaredLogger
}

func (h SolutionHandler) SendSolution(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	try := try.Try{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&try, r.PostForm)
	if err != nil {
		http.Error(w, `send solution bad form`, http.StatusBadRequest)
		return
	}
	try.Created = time.Now()

	err = h.Transmitter.Transmit("solution", try)
	if err != nil {
		h.Logger.Error("solution transmit failed", zap.Error(err))
		http.Error(w, `{"error": "publish to broker failed"}`, http.StatusInternalServerError)
		return
	}

	response.WriteResponse(w, http.StatusOK, true, "success")
}
