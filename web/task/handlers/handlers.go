package handlers

import (
	"fmt"
	"github.com/artkescha/checker/online_checker/pkg/fileStorage"
	"github.com/artkescha/checker/online_checker/pkg/kit"
	"github.com/artkescha/checker/online_checker/pkg/session"
	"github.com/artkescha/checker/online_checker/pkg/task"
	"github.com/artkescha/checker/online_checker/pkg/task/repository"
	"github.com/artkescha/checker/online_checker/pkg/unzipper"
	"github.com/artkescha/checker/online_checker/pkg/zipper"
	"github.com/artkescha/checker/online_checker/web/request"
	"strconv"

	"github.com/artkescha/checker/online_checker/web/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"path/filepath"
)

type Tasker interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	EditForm(w http.ResponseWriter, r *http.Request)
	ReadOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	CreateForm(w http.ResponseWriter, r *http.Request)

	UploadTests(w http.ResponseWriter, r *http.Request)
	DownloadTests(w http.ResponseWriter, r *http.Request)
}

type TaskHandler struct {
	Tmpl           *template.Template
	TasksRepo      repository.TaskRepo
	SessionManager session.Manager
	Logger         *zap.SugaredLogger
}

func (h TaskHandler) List(w http.ResponseWriter, r *http.Request) {
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
		Tasks []task.Task
		UserId int64
	}{
		Tasks: tasks,
		UserId:user_.ID,
	})
	if err != nil {
		h.Logger.Error("tasks list executeTemplate err", err)
		http.Error(w, fmt.Sprintf(`tasks list template err %s`, err), http.StatusInternalServerError)
		return
	}
}

func (h TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	/*err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("decode request params failed %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("request: %v", r)
	task := task.Task{}
	fmt.Printf("task : %v", task)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)



	//err = decoder.Decode(&task, r.PostForm)
	//if err != nil {
	//	http.Error(w, `Bad form`, http.StatusBadRequest)
	//	return
	//}
    */

	task := task.Task{}

	err := request.DecodePostParams(&task, r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//TODO!!!!!!!!!!!
	task.TestsPath = "C:/"
	_, err = h.TasksRepo.Insert(r.Context(), task)
	if err != nil {
		h.Logger.Error("create task err", err)
		http.Error(w, fmt.Sprintf("create task err %s", err), http.StatusInternalServerError)
		return
	}
	response.WriteResponse(w, http.StatusOK, true, "success")

	//http.Redirect(w, r, "/admin", http.StatusSeeOther)
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

func (h *TaskHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	task, err := h.TasksRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "description.html", task)
	if err != nil {
		http.Error(w, `Template err`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("read task %v with id: %d", task, id)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}
	/*r.ParseForm()
	task := new(task.Task)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	fmt.Printf("task: %v", task)
	err = decoder.Decode(task, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}*/
	task := new(task.Task)

	err = request.DecodePostParams(task, r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err)
		return
	}

	task.ID = id

	ok, err := h.TasksRepo.Update(r.Context(), task)
	if err != nil {
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", task, ok)


	response.WriteResponse(w, http.StatusOK, ok, "success")
	//http.Redirect(w, r, "/admin", http.StatusSeeOther)
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

func (h TaskHandler) SolutionForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["taskID"])
	if err != nil {
		http.Error(w, `bad id`, http.StatusBadRequest)
		return
	}
	/*user, err := request.ExtractContext(r)
	if err != nil {
		h.Logger.Error("extract use request context err", err)
		http.Error(w, `extract use request context err`, http.StatusUnauthorized)
		return
	}*/
	//TODO middelware later
	//TODO replace later languageID = 1 (golang 1.13)
	err = h.Tmpl.ExecuteTemplate(w, "send_solution.html", struct {
		TaskID     int
		LanguageID int
		UserID int64
	}{TaskID: taskID, LanguageID: 1, UserID:64})

	if err != nil {
		h.Logger.Error("execute send solution template err", err)
		http.Error(w, `send solution template err`, http.StatusInternalServerError)
		return
	}
}

func (h TaskHandler) UploadTests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["taskID"]
	if taskID == "" {
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("upload file failed, reason: task id is empty"))
		return
	}
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myArch`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	var uploadError error

	file, handler, err := r.FormFile("myArch")
	if err != nil {
		h.Logger.Errorf("get file failed, reason %s", uploadError)
		response.WriteError(w, http.StatusInternalServerError, fmt.Errorf("uploaded archive failed,reason: %s", err))
		return
	}

	defer func() {
		if uploadError != nil {
			response.WriteError(w, http.StatusInternalServerError, fmt.Errorf("uploaded archive %s failed,reason: %s", handler.Filename, uploadError))
			return
		}
		response.WriteResponse(w, http.StatusOK, fmt.Sprintf("archive %s uploaded successfully", handler.Filename))
	}()

	h.Logger.Debugf("Uploaded File: %+v\n", handler.Filename)
	h.Logger.Debugf("File Size: %+v\n", handler.Size)
	h.Logger.Debugf("MIME Header: %+v\n", handler.Header)

	rootPath := "/home/artyom/solutions/tests"
	tempPath := "/home/artyom/temp-zip"

	//make rootPath
	fStorage, err := fileStorage.New(tempPath, rootPath)

	if err != nil {
		uploadError = err
		h.Logger.Errorf("upload file failed, reason %s", uploadError)
		return
	}

	tempZipPath, err := fStorage.UploadFile(file)
	//required!!
	file.Close()

	if err != nil {
		uploadError = err
		h.Logger.Errorf("upload file failed, reason %s", uploadError)
		return
	}

	unzpr := unzipper.New()
	path := filepath.Join(rootPath, taskID)

	if exist := kit.ExistsFolder(path); !exist {
		if err := kit.EnsureDir(path); err != nil {
			uploadError = err
			h.Logger.Errorf("upload file failed, create path %s failed %s", taskID, err)
			return
		}
	}
	if err := unzpr.Unzip(tempZipPath, path, []string{".in", ".out"}); err != nil {
		uploadError = err
		h.Logger.Errorf("upload file failed, unzip archive failed %s", err)
	}
	h.Logger.Info(tempZipPath)
	h.Logger.Info(path)
}

func (h TaskHandler) DownloadTests(w http.ResponseWriter, r *http.Request) {
	var downloadError error

	vars := mux.Vars(r)
	taskID := vars["taskID"]
	if taskID == "" {
		h.Logger.Error("task id is empty")
		response.WriteError(w, http.StatusBadRequest, fmt.Errorf("upload file failed, reason: task id is empty"))
		return
	}

	defer func() {
		if downloadError != nil {
			response.WriteError(w, http.StatusBadRequest, fmt.Errorf("failedly download file %s\n", downloadError))
			return
		}
		response.WriteResponse(w, http.StatusOK, "archive download successfully")
	}()

	rootPath := filepath.Join("../tests/", taskID)

	files, err := kit.FilePathWalkDir(rootPath, []string{".in", ".out"})
	if err != nil {
		downloadError = err
		h.Logger.Errorf("error scanning files failed: %s", err)
		return
	}

	rootArchPath := "../temp-zip"

	zpr, err := zipper.New(rootArchPath, ".zip")
	if err != nil {
		downloadError = err
		h.Logger.Errorf("new zipper with prefix .zip failed %s", err)
		return
	}

	err = zpr.Add(files, taskID)

	if err != nil {
		downloadError = err
		h.Logger.Errorf("add file to zip archive %s", err)
		return
	}
	archName, err := zpr.Get()
	if err != nil {
		downloadError = err
		h.Logger.Errorf("get archive failed %s\n", err)
		return
	}

	//make rootPath
	fStorage, err := fileStorage.New("../temp-zip", rootArchPath)
	if err != nil {
		downloadError = err
		h.Logger.Errorf("create file storage failed %s: ", err)
		return
	}
	fileBytes, err := fStorage.DownloadFile(archName)
	if err != nil {
		downloadError = err
		h.Logger.Errorf("file storage download file with name %s failed %s: ", archName, err)
		return
	}
	w.Write(fileBytes)
}

func (h TaskHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		h.Logger.Error("create task executeTemplate err", err)
		http.Error(w, `create task template err`, http.StatusInternalServerError)
		return
	}
}
