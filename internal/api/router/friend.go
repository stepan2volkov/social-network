package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stepan2volkov/social-network/internal/api/openapi/friendopenapi"
	"github.com/stepan2volkov/social-network/internal/entities/token"
)

var _ friendopenapi.ServerInterface = &Router{}

// FollowUser implements friendopenapi.ServerInterface
func (rt *Router) FollowUser(w http.ResponseWriter, r *http.Request) {
	u := friendopenapi.FollowRequest{}
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentUser, err := token.GetProfileFromCtx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.profileApp.Subscribe(r.Context(), currentUser.Username, u.LeaderUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetUserFriends implements friendopenapi.ServerInterface
func (rt *Router) GetUserFriends(w http.ResponseWriter, r *http.Request, username string) {
	friends, err := rt.profileApp.GetLeaders(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]friendopenapi.Friend, 0, len(friends))

	for _, f := range friends {
		resp = append(resp, friendopenapi.Friend{
			Birthdate: types.Date{
				Time: f.Birthdate,
			},
			CityId:    f.CityID,
			CreatedAt: f.CreatedAt,
			Firstname: f.Firstname,
			Gender:    string(f.Gender),
			Id:        uint64(f.ID),
			Lastname:  f.Lastname,
			Status:    string(f.Status),
			Username:  username,
		})
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
