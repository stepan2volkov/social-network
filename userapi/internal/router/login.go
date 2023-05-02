package router

import (
	"errors"
	"net/http"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

// Login implements openapi.ServerInterface.
func (rt *Router) Login(w http.ResponseWriter, r *http.Request) {
	req, err := renderRequest[userapi.LoginJSONRequestBody](r.Body)
	defer r.Body.Close()
	if err != nil {
		logError("Login", "error to parse request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := rt.userApp.Login(r.Context(), domain.Credentials{
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, domain.ErrUnauthorized) {
		logError("Login", "unauthorized", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		logError("Login", "error to get user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = rt.sessionManager.RenewToken(r.Context()); err != nil {
		logError("Login", "error to renew session token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rt.sessionManager.Put(r.Context(), "user", user)
}

func (rt *Router) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		needAuth := ctx.Value(userapi.CookieAuthScopes) != nil

		if needAuth {
			if !rt.sessionManager.Exists(ctx, "user") {
				// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				w.WriteHeader(http.StatusUnauthorized)

				return
			}
			ctx = domain.SetUserIdentity(ctx, rt.sessionManager.Get(ctx, "user").(domain.UserIdentity))
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
