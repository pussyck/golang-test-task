package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func LoadDataHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed read file")
		return
	}

	err = processJSON(data)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to process JSON data")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "Data loaded successfully"}`)
}

func writeErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := ErrorResponse{
		Code:    code,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		http.Error(w, "Failed to encode error response", 500)
	}
}
