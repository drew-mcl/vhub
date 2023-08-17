package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RespondWithJSON writes JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// ParseJSONRequest parses JSON from the request body and decodes it into the given struct
func ParseJSONRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

// GetParam returns a URL parameter from the request
func GetParam(r *http.Request, key string) string {
	return mux.Vars(r)[key]
}
