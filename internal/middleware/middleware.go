package midlleware

import (
	"app_burse_backend/internal/app"
	"net/http"
)

type Middlewares struct {
	app app.AppContext
}

func NewMiddlewares(app app.AppContext) *Middlewares {
	return &Middlewares{app: app}
}

func (mw *Middlewares) NotLoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		if !mw.isLoggedIn(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (mw *Middlewares) isLoggedIn(r *http.Request) bool {
	token := r.Header.Get("Authorization")

	return true
}
