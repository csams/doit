package routes

import (
	"net/http"
	"strconv"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type TaskList struct {
	Length int         `json:"length"`
	Tasks  []apis.Task `json:"tasks"`
}

type TaskController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewTaskController(db *gorm.DB, log logr.Logger) *TaskController {
	return &TaskController{
		DB:  db,
		Log: log,
	}
}

// List returns the list of tasks for the user that authenticated for the
// current request
func (c *TaskController) List(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userId, err := strconv.Atoi(chi.URLParam(r, "userid"))

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	if u.ID != uint(userId) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// TODO: add pagination. this probably can be done generically in some
	// middleware in which we call db.Limit and store the resulting db object
	// on in the request context
	db := c.DB

	var results []apis.Task

	if err := db.Where("owner_id = ?", u.ID).Find(&results).Error; err != nil {
		http.Error(w, "error retrieving tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := TaskList{
		Length: len(results),
		Tasks:  results,
	}

	render.JSON(w, r, result)
}

func (c *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userId, err := strconv.Atoi(chi.URLParam(r, "userid"))

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	if u.ID != uint(userId) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	task := &apis.Task{}

	err = render.Bind(r, task)
	if err != nil {
		http.Error(w, "Unable to decode task: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !apis.IsValidStatus(task.Status) {
		http.Error(w, "task status is invalid", http.StatusBadRequest)
		return
	}

	task.OwnerId = u.ID
	task.AssigneeId = u.ID
	task.State = apis.Open

	if err = c.DB.Create(task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, task)
}

func (c *TaskController) Get(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(chi.URLParam(r, "userid"))

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	if u.ID != uint(userId) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	taskId, err := strconv.Atoi(chi.URLParam(r, "taskid"))

	if err != nil {
		http.Error(w, "invalid taskid", http.StatusBadRequest)
		return
	}

	task := &apis.Task{ID: uint(taskId), OwnerId: uint(userId)}
	if err := c.DB.First(task).Error; err != nil {
		http.Error(w, "Unable to retrieve task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, task)
}

func (c *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(chi.URLParam(r, "userid"))

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	if u.ID != uint(userId) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	taskId, err := strconv.Atoi(chi.URLParam(r, "taskid"))

	if err != nil {
		http.Error(w, "invalid taskid", http.StatusBadRequest)
		return
	}

	task := &apis.Task{ID: uint(taskId), OwnerId: uint(userId)}
	if err := c.DB.Find(task).Error; err != nil {
		http.Error(w, "Unable to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.Bind(r, task)
	task.ID = uint(taskId)
	task.OwnerId = uint(userId)

	if err := c.DB.Save(task).Error; err != nil {
		http.Error(w, "Unable to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, task)
}

func (c *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(chi.URLParam(r, "userid"))

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	if u.ID != uint(userId) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	taskId, err := strconv.Atoi(chi.URLParam(r, "taskid"))

	if err != nil {
		http.Error(w, "invalid taskid", http.StatusBadRequest)
		return
	}

	task := &apis.Task{ID: uint(taskId), OwnerId: uint(userId)}
	if err := c.DB.Find(task).Error; err != nil {
		http.Error(w, "Unable to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.DB.Delete(task).Error; err != nil {
		http.Error(w, "Unable to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, task)
}
