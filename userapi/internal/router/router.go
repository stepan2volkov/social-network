package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
	"github.com/stepan2volkov/social-network/profile/internal/app/userapp"
)

var _ userapi.ServerInterface = (*Router)(nil)

type Metrics interface {
	WithLatency(method, handler, status string, started time.Time)
}

type Router struct {
	http.Handler
	sessionManager *scs.SessionManager
	userApp        *userapp.App
	metrics        Metrics
}

func New(
	sessionManager *scs.SessionManager,
	userApp *userapp.App,
	metrics Metrics,
) *Router {
	mux := chi.NewRouter()

	rt := &Router{
		Handler:        mux,
		sessionManager: sessionManager,
		userApp:        userApp,
		metrics:        metrics,
	}

	// Middlewares.
	options := userapi.ChiServerOptions{
		BaseRouter: mux,
		Middlewares: []userapi.MiddlewareFunc{
			rt.auth,
			sessionManager.LoadAndSave,
			rt.latency,
			middleware.Recoverer,
		},
	}

	// Register routes.
	userapi.HandlerWithOptions(rt, options)

	// Metrics handler.

	mux.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Static assets.

	return rt
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}
func (rt *Router) latency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		method := r.Method
		handler := chi.RouteContext(r.Context()).RoutePattern()

		rec := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		status := strconv.Itoa((rec.statusCode / 100) * 100)
		rt.metrics.WithLatency(method, handler, status, startedAt)
	})
}
