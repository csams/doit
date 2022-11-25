package routes

import (
	"net/http"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type CommentController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewCommentController(db *gorm.DB, log logr.Logger) *CommentController {
	return &CommentController{
		DB:  db,
		Log: log,
	}
}

func (c *CommentController) List(w http.ResponseWriter, r *http.Request) {
}

func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *CommentController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *CommentController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
}
