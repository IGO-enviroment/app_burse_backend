package v1

import (
	"app_burse_backend/internal/app"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, appContext app.AppContext) {

	v1 := router.PathPrefix("/v1").Subrouter()

	// ... setup v1 routes here
	v1.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// Handle user requests
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User List"))
	}).Methods("GET")

	v1.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Handle user requests by ID
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Details"))
	}).Methods("GET")

	v1.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		// Handle user requests by current user
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Profile"))
	}).Methods("GET")
}
