package handler

import (
	"app/internal/storage"
	"app/internal/utils"
	"net/http"
)

// LoadDataHandler handler for data load route
func LoadDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}

	file, _, fileErr := r.FormFile("file")
	url := r.FormValue("url")

	if fileErr == nil {
		defer file.Close()
		storage.ProcessFile(w, file)
		return
	}

	if url != "" {
		storage.ProcessURL(w, url)
		return
	}

	utils.WriteResponse(w, http.StatusBadRequest, "No file or URL provided")
}
