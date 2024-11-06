package middlleware_test

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app/web"
	middlleware "app_burse_backend/internal/middleware"
	tokenservice "app_burse_backend/pkg/token_service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotLoggedIn(t *testing.T) {
	pwd := "../../"
	cfg := configs.NewCondfig().LoadForTest(pwd)
	app := web.NewWebContext(cfg)

	t.Run("Вернет 401 при отсутствии токена", func(t *testing.T) {
		mw := middlleware.NewMiddleware(app)
		next := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}

		handler := mw.NotLoggedInMiddleware(next)
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("Продолжит обработку запроса при валидном токене", func(t *testing.T) {
		mw := middlleware.NewMiddleware(app)
		next := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}

		handler := mw.NotLoggedInMiddleware(next)
		req := httptest.NewRequest("GET", "/api/v1/users", nil)

		validToken, err := tokenservice.NewTokenService("secret", 60).GenerateToken(10)
		if err != nil {
			t.Fatalf("Error generating token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+validToken)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}
