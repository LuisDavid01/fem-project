package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LuisDavid01/femProject/internal/store"
	"github.com/LuisDavid01/femProject/internal/tokens"
	"github.com/LuisDavid01/femProject/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}
type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		panic("user not found in context")
	}
	return user

}

func (m *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid Authorization header format"})
			return
		}

		token := headerParts[1]
		user, err := m.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid token"})
			return
		}
		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Token expired or invalid"})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})
}

func (m *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil || user.ID == 0 {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
