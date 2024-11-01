package v1

import (
	"app_burse_backend/internal/app"
	"net/http"

	"github.com/gorilla/mux"
)

type Delivery struct {
	appContext app.AppContext
}

func SetupRoutes(router *mux.Router, appContext app.AppContext) {
	users := router.PathPrefix("/v1").Subrouter().PathPrefix("/users").Subrouter()
	delivery := &Delivery{appContext: appContext}

	users.HandleFunc("/me", delivery.GetMe).Methods("GET")
	users.HandleFunc("/login", delivery.Login).Methods("POST")
}

func (d *Delivery) GetMe(w http.ResponseWriter, r *http.Request) {
	// Handle user requests
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User List"))
}

func (d *Delivery) Login(w http.ResponseWriter, r *http.Request) {
	// Handle user login requests
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Login"))
}
