package utils

import (
	"encoding/json"
	"net/http"
)

// WriteResponse write system response
func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": message}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// IsZip check file is ZIP
func IsZip(data []byte) bool {
	return len(data) > 4 && string(data[:2]) == "PK" && string(data[2:4]) == "\x03\x04"
}
