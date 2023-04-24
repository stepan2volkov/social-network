package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/stepan2volkov/social-network/profile/internal/app"
)

type Router struct {
	http.Handler
	app *app.App
}

func New(app *app.App) *Router {
	mux := chi.NewRouter()

	// Middlewares.
	mux.Use(middleware.Recoverer)

	// Register routes.

	// Static assets.

	return &Router{
		Handler: mux,
		app:     app,
	}
}
