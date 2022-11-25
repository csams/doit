package routes

import (
	"net/http"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type UserController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewUserController(db *gorm.DB, log logr.Logger) *UserController {
	return &UserController{
		DB:  db,
		Log: log,
	}
}

func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
}
