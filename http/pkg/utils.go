package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse is a utility function to write JSON responses
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ParseJSONBody parses JSON from request body into the given interface
func ParseJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON payload",
		})
		return err
	}
	return nil
}
