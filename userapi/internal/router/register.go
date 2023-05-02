package router

import (
	"net/http"

	"github.com/stepan2volkov/social-network/profile/internal/api/userapi"
	"github.com/stepan2volkov/social-network/profile/internal/domain"
)

// Register implements openapi.ServerInterface.
func (rt *Router) Register(w http.ResponseWriter, r *http.Request) {
	req, err := renderRequest[userapi.RegisterJSONRequestBody](r.Body)
	defer r.Body.Close()
	if err != nil {
		logError("Register", "error to decode request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := rt.userApp.Register(r.Context(), domain.Credentials{
		Username: req.Username,
		Password: req.Password,
	}, domain.Profile{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Birthdate: req.Birthdate.Time,
		Biography: req.Biography,
		City:      req.City,
	})
	if err != nil {
		logError("Register", "error to register user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := userapi.Register201JSONResponse{
		ID: userID,
	}

	writeResponse("Register", w, resp)
}
