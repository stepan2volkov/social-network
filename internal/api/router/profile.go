package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stepan2volkov/social-network/internal/api/openapi/profileopenapi"
)

var _ profileopenapi.ServerInterface = &Router{}

// GetUserProfile implements openapi.ServerInterface
func (rt *Router) GetUserProfile(w http.ResponseWriter, r *http.Request, username string) {
	p, err := rt.profileApp.GetProfileByUsername(r.Context(), username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(profileopenapi.UserProfile{
		Id:        uint64(p.ID),
		Username:  p.Username,
		Birthdate: types.Date{Time: p.Birthdate},
		CityId:    p.CityID,
		CreatedAt: p.CreatedAt,
		Firstname: p.Firstname,
		Gender:    string(p.Gender),
		Lastname:  p.Lastname,
	}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
