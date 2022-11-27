package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/csams/doit/pkg/apis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type UserContextKey string

const (
	UserKey UserContextKey = "userCtxKey"
)

type UserRequest struct {
	*apis.User
}

func (u *UserRequest) Bind(r *http.Request) error {
	return nil
}

type UserResponse struct {
	*apis.User
}

func (u *UserResponse) Bind(r *http.Request) error {
	return nil
}

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

func WithUser(ctx context.Context, user *UserRequest) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func UserFromContext(ctx context.Context) (*UserRequest, error) {
	obj := ctx.Value(UserKey)
	if obj == nil {
		return nil, errors.New("Expected user in request context")
	}
	req, ok := obj.(*UserRequest)
	if !ok {
		return nil, errors.New("Object stored in request context couldn't convert to *apis.User")
	}
	return req, nil
}

func (c *UserController) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user apis.User
		var err error

		// TODO: check if the requested user is the logged in user or is
		// allowed to view / update the requested user. This might go into a
		// separate middleware.
		if userID := chi.URLParam(r, "userid"); userID != "" {
			err = c.DB.First(&user, "id = ?", userID).Error
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		userReq := UserRequest{
			&user,
		}

		ctx := WithUser(r.Context(), &userReq)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
	user, _ := UserFromContext(r.Context())
	render.JSON(w, r, user)
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	req := UserRequest{}
	render.Bind(r, &req)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
}
