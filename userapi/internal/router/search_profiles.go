package router

import (
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
)

// SearchProfiles implements userapi.ServerInterface.
func (rt *Router) SearchProfiles(w http.ResponseWriter, r *http.Request, params userapi.SearchProfilesParams) {
	ctx := r.Context()

	profiles, err := rt.userApp.SearchProfiles(ctx, params.Firstname, params.Lastname)
	if err != nil {
		logError("SearchProfiles", "error to search profiles", err)
	}

	resp := make(userapi.SearchProfiles200JSONResponse, 0, len(profiles))

	for _, profile := range profiles {
		resp = append(resp, userapi.Profile{
			Firstname: profile.Firstname,
			Lastname:  profile.Lastname,
			Birthdate: types.Date{
				Time: profile.Birthdate,
			},
			Biography: profile.Biography,
			City:      profile.City,
		})
	}

	writeResponse("SearchProfiles", w, resp)
}
