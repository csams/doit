/*
Here we setup the REST handler for the service using chi.
*/
package routes

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/csams/doit/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-logr/logr"
)

// NewHandler sets up all of the routes for the site
func NewHandler(db *gorm.DB, authProvider *auth.TokenProvider, log logr.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(auth.Authenticator(db, authProvider))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	meController := NewMeController(db, log.WithName("meController"))
	userController := NewUserController(db, log.WithName("userController"))
	taskController := NewTaskController(db, log.WithName("taskController"))
	commentController := NewCommentController(db, log.WithName("commentController"))
	annotationController := NewAnnotationController(db, log.WithName("annotationController"))
	policyController := NewPolicyController(db, log.WithName("policyController"))

	r.Route("/me", func(r chi.Router) {
		r.Get("/", meController.Get)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
		r.Route("/{userid}", func(r chi.Router) {
			r.Use(userController.UserCtx)
			r.Get("/", userController.Get)
			r.Put("/", userController.Update)
			r.Delete("/", userController.Delete)

			r.Route("/shares", func(r chi.Router) {
				r.Get("/with", policyController.ListSharedWith)
				r.Get("/from", policyController.ListSharedFrom)
				r.Post("/", policyController.Create)
				r.Route("/{delegateid}", func(r chi.Router) {
					r.Get("/", policyController.Get)
					r.Put("/", taskController.Update)
					r.Delete("/", taskController.Delete)
				})
			})

			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskController.List)
				r.Post("/", taskController.Create)
				r.Route("/{taskid}", func(r chi.Router) {
					r.Get("/", taskController.Get)
					r.Put("/", taskController.Update)
					r.Delete("/", taskController.Delete)

					r.Route("/comments", func(r chi.Router) {
						r.Get("/", commentController.List)
						r.Post("/", commentController.Create)
						r.Route("/{commentid}", func(r chi.Router) {
							r.Get("/", commentController.Get)
							r.Put("/", commentController.Update)
							r.Delete("/", commentController.Delete)
						})
					})

					r.Route("/annotations", func(r chi.Router) {
						r.Get("/", annotationController.List)
						r.Post("/", annotationController.Create)
						r.Route("/{annotationid}", func(r chi.Router) {
							r.Get("/", annotationController.Get)
							r.Put("/", annotationController.Update)
							r.Delete("/", annotationController.Delete)
						})
					})
				})
			})
		})
	})
	return r
}
