package web

import (
	v1 "app_burse_backend/internal/users/delivery/http/v1"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Routes struct {
	App *WebContext

	Route *mux.Router
}

func NewRoutes(App *WebContext) *Routes {
	return &Routes{App: App}
}

func (r *Routes) Setup() *mux.Router {
	// ... setup routes here
	r.Route = mux.NewRouter()

	r.App.log.Info("dsf", zap.Bool("boolField", r.App.DB() == nil))
	v1.SetupRoutes(r.Route, r.App)

	return r.Route
}
