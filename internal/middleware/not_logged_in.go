package middlleware

import (
	tokenservice "app_burse_backend/pkg/token_service"
	"net/http"
)

func (mw *Middleware) NotLoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		if !mw.isLoggedIn(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Проверка авторизации пользователя.
// Возвращает true, если токен валидный, иначе false.
func (mw *Middleware) isLoggedIn(r *http.Request) bool {
	cfg := mw.app.Configs().Web

	service := tokenservice.NewTokenService(cfg.TokenSecret, cfg.TokenExpiration)
	token, err := service.TokenFromRequest(r)
	if err != nil || token == "" {
		return false
	}

	_, ok := service.VerifyToken(token)
	return ok
}
