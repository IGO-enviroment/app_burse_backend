package v1_test

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app/web"
	"app_burse_backend/internal/domain"
	v1 "app_burse_backend/internal/users/delivery/http/v1"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	app := web.NewWebContext(configs.NewCondfig().Load())
	app.InitDB()

	connect, err := app.DB().BeginTx(context.TODO(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defer func() {
		connect.Rollback()
		app.DB().Close()
	}()

	mockRequest := func(request any) *http.Response {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(request)
		req, err := http.NewRequest("POST", "/v1/users/login", &buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router := mux.NewRouter()
		v1.SetupRoutes(router, nil)
		router.ServeHTTP(w, req)

		return w.Result()
	}

	t.Run("Should return 200 OK", func(t *testing.T) {
		request := map[string]string{"email": "test@example.com", "password": "test"}
		user := &domain.User{Email: request["email"]}
		user.SetPassword(request["password"])
		connect.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.DigestPassword)

		res := mockRequest(request)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
