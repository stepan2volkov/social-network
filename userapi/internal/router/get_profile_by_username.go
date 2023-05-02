package router

import (
	"errors"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

// GetProfileByUsername implements authapi.ServerInterface.
func (rt *Router) GetProfileByUsername(w http.ResponseWriter, r *http.Request, username string) {
	ctx := r.Context()

	profile, err := rt.userApp.GetProfileByUsername(ctx, username)
	if errors.Is(err, domain.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		logError("GetProfileByUsername", "error to get profile from repo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := userapi.GetProfileByUsername200JSONResponse{
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		Birthdate: types.Date{Time: profile.Birthdate},
		Biography: profile.Biography,
		City:      profile.City,
	}

	writeResponse("GetProfileByUsername", w, resp)
}
