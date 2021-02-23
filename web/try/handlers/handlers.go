package handlers

import (
	"github.com/artkescha/grader/online_checker/pkg/session"
	"github.com/artkescha/grader/online_checker/pkg/tries"
	"github.com/artkescha/grader/online_checker/pkg/tries/repository"
	"github.com/artkescha/grader/online_checker/pkg/tries/transmitter"
	"github.com/artkescha/grader/online_checker/web/request"
	"github.com/artkescha/grader/online_checker/web/response"
	"log"
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
	var try try.Try
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
	try.UserID = user.ID

	log.Printf("try: %v", try)

	err = h.Transmitter.Transmit("solution", try)
	if err != nil {
		h.Logger.Error("solution transmit failed", zap.Error(err))
		http.Error(w, `{"error": "publish to broker failed"}`, http.StatusInternalServerError)
		return
	}

	response.WriteResponse(w, http.StatusOK, true, "success")
}
