package v1

import (
	"app_burse_backend/internal/app"
	midlleware "app_burse_backend/internal/middleware"
	users_entity "app_burse_backend/internal/users/entity"
	users_repository "app_burse_backend/internal/users/repo"
	users_usecase "app_burse_backend/internal/users/usecase"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Delivery struct {
	appContext app.AppContext
}

func SetupRoutes(router *mux.Router, appContext app.AppContext) {
	users := router.PathPrefix("/v1").Subrouter().PathPrefix("/users").Subrouter()
	delivery := &Delivery{appContext: appContext}

	middlware := midlleware.NewMiddlewares(appContext)

	users.HandleFunc(
		"/me",
		delivery.GetMe,
	).Methods("GET")
	users.HandleFunc(
		"/login",
		middlware.NotLoggedInMiddleware(delivery.Login),
	).Methods("POST")
}

func (d *Delivery) GetMe(w http.ResponseWriter, r *http.Request) {
	// Handle user requests
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User List"))
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *Delivery) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return
	}

	usecase := users_usecase.NewLoginUsecase(
		*users_repository.NewRepository(d.appContext.DB()),
		users_entity.LoginEntity{Email: request.Email, Password: request.Password},
		d.appContext,
	)
	result := usecase.Call()

	if !result.Success() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(result.Error()))
		return
	}

	d.appContext.Logger().Info(
		"[users.login] Пользователь авторизован.",
		zap.String("email", request.Email),
		zap.String("ip", r.RemoteAddr),
	)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result.Data())
}
