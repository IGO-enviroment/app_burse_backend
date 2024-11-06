package middlleware

import (
	"fmt"
	"net/http"
)

// Проверка, что пользователь не авторизован.
func (mw *Middleware) NotLoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		if _, ok := mw.isLoggedIn(r); ok {
			fmt.Println("User is already logged in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
