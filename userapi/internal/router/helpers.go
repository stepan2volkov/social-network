package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func renderRequest[T any](r io.Reader) (T, error) {
	var req T

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return req, err
	}

	return req, nil
}

func writeResponse[T any](routePath string, w http.ResponseWriter, resp T) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(routePath, "error to encode response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func logError(routePath, msg string, err error) {
	log.Printf("[router.%s] %s: %s\n", routePath, msg, err)
}
