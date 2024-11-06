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
	"net/http"
	"net/http/httptest"
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
	mockApp.On("Logger").Return(app.Logger())
	mockApp.On("Producer").Return(app.Producer())
	mockApp.On("Configs").Return(app.Configs())
	mockApp.On("Locales").Return(app.Locales())

	return connect, app, mockApp
}

func TestLogin(t *testing.T) {
	connect, app, mockApp := setup(t)

	defer func() {
		connect.Rollback()
		app.DB().(*sqlx.DB).Close()
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
		v1.SetupRoutes(router, mockApp)
		router.ServeHTTP(w, req)

		return w.Result()
	}

	t.Run("Should return 200 OK", func(t *testing.T) {
		request := map[string]string{"email": "test@example.com", "password": "test"}
		user := &domain.User{Email: request["email"]}
		user.SetPassword(request["password"])
		_, err := users_repository.NewRepository(mockApp.DB()).Create(
			[]service.FieldDB{
				{Name: "email", Value: user.Email},
				{Name: "digest_password", Value: user.DigestPassword},
			},
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		res := mockRequest(request)

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
