package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/stepan2volkov/social-network/internal/api/openapi/authopenapi"
	"github.com/stepan2volkov/social-network/internal/api/openapi/friendopenapi"
	"github.com/stepan2volkov/social-network/internal/api/openapi/profileopenapi"
	"github.com/stepan2volkov/social-network/internal/app/authapp"
	"github.com/stepan2volkov/social-network/internal/app/profileapp"
	"github.com/stepan2volkov/social-network/internal/entities/token"
)

type VersionInfo struct {
	BuildCommit string
	BuildTime   string
}

type Router struct {
	http.Handler
	version    VersionInfo
	authApp    *authapp.App
	profileApp *profileapp.App
}

func New(
	version VersionInfo,
	authApp *authapp.App,
	profileApp *profileapp.App,
) *Router {
	r := chi.NewRouter()

	rt := &Router{
		version:    version,
		authApp:    authApp,
		profileApp: profileApp,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/auth", authopenapi.Handler(rt))
	r.Mount("/profiles", profileopenapi.HandlerWithOptions(rt, profileopenapi.ChiServerOptions{
		Middlewares: []profileopenapi.MiddlewareFunc{rt.authMiddleware},
	}))
	r.Mount("/friends", friendopenapi.HandlerWithOptions(rt, friendopenapi.ChiServerOptions{
		Middlewares: []friendopenapi.MiddlewareFunc{rt.authMiddleware},
	}))

	r.Get("/__version__", rt.versionHandler)
	r.Get("/__heartbeat__", rt.heartbeatHandler)

	rt.Handler = r

	return rt
}

func (rt *Router) versionHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(&rt.version); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rt *Router) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (rt *Router) authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Access-Token")

		if !strings.HasPrefix(accessToken, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		accessToken = accessToken[len("Bearer "):]

		userProfile, err := rt.authApp.Verify(r.Context(), accessToken)
		if err != nil {
			log.Printf("Auth Middleware: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := token.CreateCtxWithProfile(r.Context(), userProfile)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
