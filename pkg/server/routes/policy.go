package routes

import (
	"net/http"
	"strconv"

	"github.com/csams/doit/pkg/auth"
	"github.com/go-chi/chi/v5"
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
}

func (c *PolicyController) ListSharedFrom(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Get(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Create(w http.ResponseWriter, r *http.Request) {
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
}

func (c *PolicyController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *PolicyController) Delete(w http.ResponseWriter, r *http.Request) {
}
