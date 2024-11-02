package middlleware

import (
	"app_burse_backend/internal/app"
	tokenservice "app_burse_backend/pkg/token_service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	app app.AppContext
}

func NewMiddleware(app app.AppContext) *Middleware {
	return &Middleware{app: app}
}

// Проверка авторизации пользователя.
// Возвращает true, если токен валидный, иначе false.
func (mw *Middleware) isLoggedIn(r *http.Request) (*jwt.Token, bool) {
	cfg := mw.app.Configs().Web

	service := tokenservice.NewTokenService(cfg.TokenSecret, cfg.TokenExpiration)
	token, err := service.TokenFromRequest(r)
	if err != nil || token == "" {
		return nil, false
	}

	return service.VerifyToken(token)
}
