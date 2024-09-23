package handler

import (
	"app/internal/storage"
	"app/internal/utils"
	"net/http"
)

// LoadDataHandler handler for load data from http or file
func (h *Handler) LoadDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}

	file, _, fileErr := r.FormFile("file")
	url := r.FormValue("url")

	if fileErr == nil {
		defer file.Close()
		err := storage.ProcessFile(file, h.RedisClient)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusOK, "data loaded successfully")
		return
	}

	if url != "" {
		err := storage.ProcessURL(url, h.RedisClient)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusOK, "data loaded successfully")
		return
	}

	utils.WriteResponse(w, http.StatusBadRequest, "No file or URL provided")
}
