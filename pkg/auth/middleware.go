package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/csams/doit/pkg/apis"
	"gorm.io/gorm"
)

type userClaims struct {
	Name     string `json:"name"`
	Username string `json:"preferred_username"`
	Audience string `json:"aud"`
}

func Authenticator(db *gorm.DB, provider *TokenProvider, clientId string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get the token from the request
			rawToken := ""
			for _, getToken := range TokenGetters {
				rawToken = getToken(r)
				if rawToken != "" {
					break
				}
			}

			// ensure we got one
			if rawToken == "" {
				http.Error(w, "No credentials supplied", http.StatusUnauthorized)
				return
			}

			// verify and parse it
			tok, err := provider.Verify(rawToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// extract the claims we care about
			u := &userClaims{}
			tok.Claims(u)
			if u.Name == "" || u.Username == "" {
				http.Error(w, "Invalid claims. Require name and preferred_username", http.StatusBadRequest)
				return
			}

			if u.Audience != clientId {
				http.Error(w, "Invalid audience.", http.StatusBadRequest)
				return
			}

			// fetch the user from the database or create them if they doesn't exist
			usr := &apis.User{
				Name:     u.Name,
				Username: u.Username,
			}
			if err := db.FirstOrCreate(usr).Error; err != nil {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// send them down the handler chain
			ctx := NewContext(r.Context(), usr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewContext(ctx context.Context, user *apis.User) context.Context {
	ctx = context.WithValue(ctx, UserCtxKey, user)
	return ctx
}

func UserFromContext(ctx context.Context) (*apis.User, error) {
	user, ok := ctx.Value(UserCtxKey).(*apis.User)
	if !ok {
		return nil, errors.New("could not retrieve user from context")
	}
	return user, nil
}

func TokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:6]) == "bearer" {
		return bearer[7:]
	}
	return ""
}
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func TokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}

type contextKey struct {
	name string
}

type TokenGetter func(*http.Request) string

var (
	UserCtxKey   = &contextKey{"token"}
	ErrorCtxKey  = &contextKey{"error"}
	TokenGetters = []TokenGetter{TokenFromHeader, TokenFromCookie, TokenFromQuery}
)
