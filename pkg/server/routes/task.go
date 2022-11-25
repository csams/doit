package routes

import (
	"net/http"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

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

func (c *TaskController) List(w http.ResponseWriter, r *http.Request) {
}

func (c *TaskController) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *TaskController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *TaskController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
}
