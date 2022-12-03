package routes

import (
	"net/http"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/auth"
	"github.com/go-logr/logr"
	"gorm.io/gorm"

	"github.com/go-chi/render"
)

type MeController struct {
	DB  *gorm.DB
	Log logr.Logger
}

func NewMeController(db *gorm.DB, log logr.Logger) *MeController {
	return &MeController{
		DB:  db,
		Log: log,
	}
}

func (c *MeController) Get(w http.ResponseWriter, r *http.Request) {
	u, err := auth.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := c.DB.Preload("AssignedTasks", "state = ?", apis.Open).First(u).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, u)
}
