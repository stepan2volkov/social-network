package router

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
	"github.com/stepan2volkov/social-network/profile/internal/app/userapp"
)

var _ userapi.ServerInterface = (*Router)(nil)

type Router struct {
	http.Handler
	sessionManager *scs.SessionManager
	userApp        *userapp.App
}

func New(
	sessionManager *scs.SessionManager,
	userApp *userapp.App,
) *Router {
	mux := chi.NewRouter()

	rt := &Router{
		Handler:        mux,
		sessionManager: sessionManager,
		userApp:        userApp,
	}

	// Middlewares.
	options := userapi.ChiServerOptions{
		BaseRouter: mux,
		Middlewares: []userapi.MiddlewareFunc{
			rt.auth,
			sessionManager.LoadAndSave,
			middleware.Recoverer,
		},
	}

	// Register routes.
	userapi.HandlerWithOptions(rt, options)

	// Static assets.

	return rt
}
