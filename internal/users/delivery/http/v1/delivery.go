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
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func SetupRoutes(router *mux.Router, appContext app.AppContext) {
	users := router.PathPrefix("/v1").Subrouter().PathPrefix("/users").Subrouter()
	delivery := NewDelivery(appContext)

	middlware := midlleware.NewMiddleware(appContext)

	users.HandleFunc(
		"/me",
		delivery.GetMe,
	).Methods("GET")
	users.HandleFunc(
		"/login",
		middlware.NotLoggedInMiddleware(delivery.Login),
	).Methods("POST")
}

type Delivery struct {
	appContext app.AppContext
}

func NewDelivery(appContext app.AppContext) *Delivery {
	return &Delivery{appContext: appContext}
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

	store := sessions.NewCookieStore([]byte(d.appContext.Configs().Web.TokenSecret))
	session, err := store.Get(r, d.appContext.Configs().Web.CookiesField)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = result.Data()
	session.Options.MaxAge = d.appContext.Configs().Web.TokenExpiration

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	d.appContext.Logger().Info(
		"[users.login] Пользователь авторизован.",
		zap.String("email", request.Email),
		zap.String("ip", r.RemoteAddr),
	)

	w.WriteHeader(http.StatusOK)
}
