/*
Here we setup the REST handler for the service using chi.
*/
package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// NewHandler sets up all of the routes for the site
func NewHandler(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", CreateUser(db))
		r.Route("/{userid}", func(r chi.Router) {
			r.Get("/", GetUser(db))
			r.Put("/", UpdateUser(db))
			r.Delete("/", DeleteUser(db))

			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", ListTasks(db))
				r.Post("/", CreateTask(db))
				r.Route("/{taskid}", func(r chi.Router) {
					r.Get("/", GetTask(db))
					r.Put("/", UpdateTask(db))
					r.Delete("/", DeleteTask(db))

					r.Route("/comments", func(r chi.Router) {
						r.Get("/", ListComments(db))
						r.Post("/", CreateComment(db))
						r.Route("/{commentid}", func(r chi.Router) {
							r.Get("/", GetComment(db))
							r.Put("/", UpdateComment(db))
							r.Delete("/", DeleteComment(db))
						})
					})

					r.Route("/annotations", func(r chi.Router) {
						r.Get("/", ListAnnotations(db))
						r.Post("/", CreateAnnotation(db))
						r.Route("/{annotationid}", func(r chi.Router) {
							r.Get("/", GetAnnotation(db))
							r.Put("/", UpdateAnnotation(db))
							r.Delete("/", DeleteAnnotation(db))
						})
					})
				})
			})
		})
	})
	return r
}
