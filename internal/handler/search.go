package handler

import (
	"app/internal/storage"
	"app/internal/utils"
	"encoding/json"
	"net/http"
)

// GetParkingDataHandler search data by query
func (h *Handler) GetParkingDataHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	globalID := query.Get("global_id")
	mode := query.Get("mode")
	id := query.Get("id")

	if globalID == "" && mode == "" && id == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "No valid search parameters provided")
		return
	}

	result, err := storage.SearchData(globalID, mode, id, h.RedisClient)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		utils.WriteResponse(w, http.StatusNotFound, "No data found")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err.Error())
	}
}
