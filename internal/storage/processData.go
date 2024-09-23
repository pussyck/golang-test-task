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

// ProcessFile saves data from a file
func ProcessFile(file io.Reader, redisClient redis.C) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file")
	}

	err = ProcessJSON(data, redisClient)
	if err != nil {
		return fmt.Errorf("failed to process JSON: %v", err)
	}

	return nil
}

// ProcessURL saves data from a URL
func ProcessURL(url string, redisClient redis.C) error {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file from URL")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read downloaded file")
	}

	if utils.IsZip(data) {
		err := ProcessZip(data, redisClient)
		if err != nil {
			return fmt.Errorf("failed to process ZIP: %v", err)
		}
	} else {
		err := ProcessJSON(data, redisClient)
		if err != nil {
			return fmt.Errorf("failed to process JSON: %v", err)
		}
	}

	return nil
}

// ProcessZip unzips and uses processJSON
func ProcessZip(data []byte, redisClient redis.C) error {
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

			err = ProcessJSON(fileData, redisClient)
			if err != nil {
				return fmt.Errorf("failed to process JSON file in ZIP: %v", err)
			}
		}
	}

	return nil
}

// ProcessJSON saves JSON records in Redis
func ProcessJSON(data []byte, redisClient redis.C) error {
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
		id := fmt.Sprintf("index:id:%.0f", record["ID"].(float64))

		recordJSON, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %v", err)
		}

		if err := redisClient.Set(globalID, recordJSON); err != nil {
			return fmt.Errorf("failed to save data to Redis: %v", err)
		}
		if err := redisClient.SAdd(mode, globalID); err != nil {
			return fmt.Errorf("failed to add Mode tag: %v", err)
		}
		if err := redisClient.SAdd(id, globalID); err != nil {
			return fmt.Errorf("failed to add ID tag: %v", err)
		}
	}

	return nil
}
