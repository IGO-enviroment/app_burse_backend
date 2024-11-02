package middlleware

import (
	"app_burse_backend/internal/app"
)

type Middleware struct {
	app app.AppContext
}

func NewMiddleware(app app.AppContext) *Middleware {
	return &Middleware{app: app}
}
