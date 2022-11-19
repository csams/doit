/*
Here we setup the REST handler for the service using chi.

This will need to be decomposed.
*/
package routes

import (
	"net/http"

	"gorm.io/gorm"
)

type RootHandler struct {
	db       *gorm.DB
	delegate http.Handler
}

// NewHandler sets up all of the routes for the site
func NewHandler(db *gorm.DB) http.Handler {
	return &RootHandler{db, nil}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
