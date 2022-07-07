package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/stepan2volkov/social-network/internal/api/openapi/authopenapi"
	"github.com/stepan2volkov/social-network/internal/entities/user"
)

var _ authopenapi.ServerInterface = &Router{}

// Login implements openapi.ServerInterface
func (rt *Router) Login(w http.ResponseWriter, r *http.Request) {
	u := authopenapi.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := rt.authApp.Login(r.Context(), u.Username, u.Password)
	if errors.Is(err, user.ErrInvalidCredentials) {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Save tokens into cookie
	w.Header().Set("Access-Token", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	w.Header().Set("Refresh-Token", fmt.Sprintf("Bearer %s", tokens.RefreshToken))
	w.WriteHeader(http.StatusOK)
}

// Register implements openapi.ServerInterface
func (rt *Router) Register(w http.ResponseWriter, r *http.Request) {
	u := authopenapi.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = rt.authApp.Register(r.Context(), user.UserProfile{
		Username:  u.Username,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Birthdate: u.Birthdate.Time,
		Gender:    user.Gender(u.Gender),
		CityID:    u.CityId,
	}, u.Password); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
