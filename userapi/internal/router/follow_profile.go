package router

import (
	"errors"
	"net/http"

	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

// FollowProfile implements authapi.ServerInterface.
func (rt *Router) FollowProfile(w http.ResponseWriter, r *http.Request, username string) {
	ctx := r.Context()

	follower, err := domain.GetUserIdentity(ctx)
	if err != nil {
		logError("FollowProfile", "error to get current user", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = rt.userApp.Follow(ctx, follower.Username, username); err != nil {
		logError("FollowProfile", "error to follow user", err)

		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrAlreadyFollow) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
