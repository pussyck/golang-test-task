package storage

import (
	"app/internal/redis"
	"app/internal/utils"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"path/filepath"
)

// ProcessFile save data from file
func ProcessFile(w http.ResponseWriter, file io.Reader) {
	data, err := io.ReadAll(file)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	err = processJSON(data)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process JSON: %v", err))
		return
	}

	utils.WriteResponse(w, http.StatusOK, "Data loaded successfully")
}

// ProcessURL save data from url
func ProcessURL(w http.ResponseWriter, url string) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		utils.WriteResponse(w, http.StatusInternalServerError, "Failed to download file from URL")
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, "Failed to read downloaded file")
		return
	}

	if utils.IsZip(data) {
		err := processZip(data)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process ZIP: %v", err))
			return
		}
	} else {
		err := processJSON(data)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process JSON: %v", err))
			return
		}
	}

	utils.WriteResponse(w, http.StatusOK, "Data loaded successfully")
}

// processZip unzip and use processJSON
func processZip(data []byte) error {
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to open ZIP archive: %v", err)
	}

	for _, file := range archive.File {
		if filepath.Ext(file.Name) == ".json" {
			rc, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in ZIP: %v", err)
			}

			fileData, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return fmt.Errorf("failed to read file in ZIP: %v", err)
			}

			err = processJSON(fileData)
			if err != nil {
				return fmt.Errorf("failed to process JSON file in ZIP: %v", err)
			}
		}
	}

	return nil
}

// processJSON save json records in redis
func processJSON(data []byte) error {
	decoder := charmap.Windows1251.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		return fmt.Errorf("failed to decode file data")
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(utf8Data, &records); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	for _, record := range records {

		globalID := fmt.Sprintf("%.0f", record["global_id"].(float64))

		mode := fmt.Sprintf("index:mode:%s", record["Mode"].(string))
		fmt.Println(globalID, mode)
		id := fmt.Sprintf("index:id:%.0f", record["ID"].(float64))

		recordJSON, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %v", err)
		}

		if err := redis.Set(globalID, recordJSON); err != nil {
			return fmt.Errorf("failed to save data to Redis: %v", err)
		}
		if err := redis.SAdd(mode, globalID); err != nil {
			return fmt.Errorf("failed to add Mode tag: %v", err)
		}
		if err := redis.SAdd(id, globalID); err != nil {
			return fmt.Errorf("failed to add ID tag: %v", err)
		}
	}

	return nil
}
