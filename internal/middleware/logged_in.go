package middlleware

import (
	users_repository "app_burse_backend/internal/users/repo"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "LoggedUser"

func (mw *Middleware) LoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int

		token, ok := mw.isLoggedIn(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok = token.Claims.(jwt.MapClaims)["id"].(int)
		if !ok {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		user, err := users_repository.NewRepository(mw.app.DB()).GetById(userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
