package routes

import (
	"net/http"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type AnnotationController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewAnnotationController(db *gorm.DB, log logr.Logger) *AnnotationController {
	return &AnnotationController{
		DB:  db,
		Log: log,
	}
}

func (c *AnnotationController) List(w http.ResponseWriter, r *http.Request) {
}

func (c *AnnotationController) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *AnnotationController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *AnnotationController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *AnnotationController) Delete(w http.ResponseWriter, r *http.Request) {
}
