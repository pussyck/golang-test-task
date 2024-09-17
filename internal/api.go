package internal

import (
	"encoding/json"
	"golang.org/x/text/encoding/charmap"
	"io"
	"mime/multipart"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func LoadDataHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Failed read file")
		return
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Failed to decode file data")
		return
	}

	err = processJSON(utf8Data)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Failed to process JSON data")
		return
	}

	writeResponse(w, http.StatusOK, "Data loaded successfully")
}

func GetParkingDataHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if r.Method != http.MethodGet {
		writeResponse(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}
	for key, values := range query {
		for _, value := range values {
			if value != "" {
				result, err := GetParkingDataByField(key, value)
				if err != nil {
					writeResponse(w, http.StatusNotFound, err.Error())
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(result); err != nil {
					http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
					return
				}
				return
			}
		}
	}
	writeResponse(w, http.StatusBadRequest, "No valid search parameters provided")
}

func writeResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "", 500)
	}
}
