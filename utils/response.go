package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse adalah format standar response API
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON menulis response JSON dengan status code tertentu
func WriteJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError helper untuk mengirim error standar
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, JSONResponse{
		Success: false,
		Message: message,
	})
}
