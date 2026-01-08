package httputil

import (
	"encoding/json"
	"net/http"
)

func JSONReponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func SuccessResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
