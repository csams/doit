package routes

import (
	"net/http"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type PolicyController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewPolicyController(db *gorm.DB, log logr.Logger) *PolicyController {
	return &PolicyController{
		DB:  db,
		Log: log,
	}
}

func (c *PolicyController) ListSharedWith(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) ListSharedFrom(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Delete(w http.ResponseWriter, r *http.Request) {
}
