package v1_test

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app"
	"app_burse_backend/internal/app/web"
	"app_burse_backend/internal/domain"
	"app_burse_backend/internal/service"
	v1 "app_burse_backend/internal/users/delivery/http/v1"
	users_repository "app_burse_backend/internal/users/repo"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/producer"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockWebContext struct {
	mock.Mock
}

func (m *MockWebContext) DB() postgres.Database {
	return m.Called().Get(0).(postgres.Database)
}

func (m *MockWebContext) Logger() *zap.Logger {
	return m.Called().Get(0).(*zap.Logger)
}

func (m *MockWebContext) Producer() *producer.Producer {
	return m.Called().Get(0).(*producer.Producer)
}

func (m *MockWebContext) Configs() *configs.Config {
	return m.Called().Get(0).(*configs.Config)
}

func (m *MockWebContext) Locales() *i18n.Localizer {
	return m.Called().Get(0).(*i18n.Localizer)
}

func setup(t *testing.T) (*sqlx.Tx, app.AppContext, app.AppContext) {
	t.Helper()

	pwd := "../../../../../"
	cfg := configs.NewCondfig().LoadForTest(pwd)
	app := web.NewWebContext(cfg)
	app.InitLocales(pwd)
	app.InitDB()

	connect, err := app.DB().(*sqlx.DB).BeginTxx(context.TODO(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockApp := new(MockWebContext)
	mockApp.On("DB").Return(connect)
	mockApp.On("Logger").Return(zap.NewNop())
	mockApp.On("Producer").Return(app.Producer())
	mockApp.On("Configs").Return(app.Configs())
	mockApp.On("Locales").Return(app.Locales())

	return connect, app, mockApp
}

func mockRequest(t *testing.T, request any, mockApp app.AppContext) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(request)
	req, err := http.NewRequest("POST", "/v1/users/login", &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	v1.SetupRoutes(router, mockApp)
	router.ServeHTTP(w, req)

	return w.Result()
}

func createUser(t *testing.T, email, password string, db postgres.Database) {
	user := &domain.User{Email: email}
	user.SetPassword(password)
	_, err := users_repository.NewRepository(db).Create(
		[]service.FieldDB{
			{Name: "email", Value: user.Email},
			{Name: "digest_password", Value: user.DigestPassword},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSuccessLogin(t *testing.T) {
	requests := []map[string]string{
		{"email": "test@example.com", "password": "test"},
		{"email": "test2@example.com", "password": "password2"},
	}

	for _, request := range requests {
		t.Run("Вернет статус 200", func(t *testing.T) {
			connect, app, mockApp := setup(t)
			t.Cleanup(func() {
				connect.Rollback()
				app.DB().(*sqlx.DB).Close()
			})

			createUser(t, request["email"], request["password"], connect)
			res := mockRequest(t, request, mockApp)

			assert.Equal(t, http.StatusOK, res.StatusCode)
		})

		t.Run("Запишет токен в заголовок", func(t *testing.T) {
			connect, app, mockApp := setup(t)
			t.Cleanup(func() {
				connect.Rollback()
				app.DB().(*sqlx.DB).Close()
			})

			createUser(t, request["email"], request["password"], connect)

			res := mockRequest(t, request, mockApp)

			fmt.Println(res.Header)
			fmt.Println("successful login", res.Header.Get("Set-Cookie"))
			fmt.Println(mockApp.Configs().Web.CookiesField, "dsfsdfsdf")
			assert.Contains(t, res.Header.Get("Set-Cookie"), mockApp.Configs().Web.CookiesField)
		})
	}
}

func TestFailedLogin(t *testing.T) {
	t.Run("Неверный пароль", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		createUser(t, "test@example.com", "test", connect)
		res := mockRequest(t, map[string]string{"email": "test@example.com", "password": "wrong"}, mockApp)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Неизвестный email", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		res := mockRequest(t, map[string]string{"email": "unknown@example.com", "password": "test"}, mockApp)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Отсутствие пароля", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		res := mockRequest(t, map[string]string{"email": "test@example.com"}, mockApp)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Отсутствие email", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		res := mockRequest(t, map[string]string{"password": "pass"}, mockApp)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Отсутствие оба поля", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		res := mockRequest(t, map[string]string{}, mockApp)

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Уже авторизованный пользователь, вернет ошибку", func(t *testing.T) {
		connect, app, mockApp := setup(t)
		t.Cleanup(func() {
			connect.Rollback()
			app.DB().(*sqlx.DB).Close()
		})

		createUser(t, "test@example.com", "test", connect)
		res := mockRequest(t, map[string]string{"email": "test@example.com", "password": "test"}, mockApp)

		fmt.Println("sdlfsdfmkdsm ", res.Header.Get("Set-Cookie"))
		cookie := &http.Cookie{
			Name:     mockApp.Configs().Web.CookiesField,
			Value:    strings.Split(res.Header.Get("Set-Cookie"), "=")[1],
			HttpOnly: true,
			MaxAge:   mockApp.Configs().Web.TokenExpiration,
		}

		req, err := http.NewRequest("POST", "/v1/users/login", nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		http.SetCookie(w, cookie)
		router := mux.NewRouter()
		v1.SetupRoutes(router, mockApp)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}
