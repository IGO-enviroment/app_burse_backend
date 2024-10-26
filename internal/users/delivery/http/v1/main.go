package v1

import (
	"app_burse_backend/internal/app"
	"app_burse_backend/pkg/queue/job"
	"app_burse_backend/pkg/queue/producer"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, appContext app.AppContext) {

	v1 := router.PathPrefix("/v1").Subrouter()

	// ... setup v1 routes here
	v1.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// Handle user requests

		appContext.Producer().Add(
			producer.AddOptions{
				DB:        appContext.DB(),
				QueueName: "default",
				RunAt:     time.Now(),
				Job: job.RawJob{
					Method: "send_email",
					Params: map[string]interface{}{
						"name":  "John Doe",
						"email": "john.doe@example.com",
					},
				},
			},
		)

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
