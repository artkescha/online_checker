package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gitlab.com/artkescha/grader/online_checker/pkg/session"
	"gitlab.com/artkescha/grader/online_checker/pkg/task"
	"gitlab.com/artkescha/grader/online_checker/pkg/task/repository"
	"gitlab.com/artkescha/grader/online_checker/web/response"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"
)

type Tasker interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	CreateForm(w http.ResponseWriter, r *http.Request)
}

type TaskHandler struct {
	Tmpl           *template.Template
	TasksRepo      repository.TaskRepo
	SessionManager session.Manager
	Logger         *zap.SugaredLogger
}

func (h TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	//TODO limit:3 offset:0 in request
	tasks, err := h.TasksRepo.List(r.Context(), 100, 0, "created_at")
	if err != nil {
		h.Logger.Error("get tasks list err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "list.html", struct {
		Tasks []task.Task
	}{
		Tasks: tasks,
	})
	if err != nil {
		h.Logger.Error("tasks list executeTemplate err", err)
		http.Error(w, `tasks list template err`, http.StatusInternalServerError)
		return
	}
}

func (h TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	task := task.Task{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&task, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	//TODO!!!!!!!!!!!
	task.TestsPaths = "C:/"
	_, err = h.TasksRepo.Insert(r.Context(), task)
	if err != nil {
		h.Logger.Error("create task err", err)
		http.Error(w, fmt.Sprintf("create task err %s", err), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h TaskHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	task, err := h.TasksRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, `no task`, http.StatusNotFound)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "edit.html", task)
	if err != nil {
		http.Error(w, `Template err`, http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}
	r.ParseForm()
	task := new(task.Task)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(task, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	task.ID = id

	ok, err := h.TasksRepo.Update(r.Context(), task)
	if err != nil {
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", task, ok)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	ok, err := h.TasksRepo.Delete(r.Context(), uint32(id))
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}
	response.WriteResponse(w, http.StatusOK, ok, "success")
}

func (h TaskHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		h.Logger.Error("create task executeTemplate err", err)
		http.Error(w, `create task template err`, http.StatusInternalServerError)
		return
	}
}
