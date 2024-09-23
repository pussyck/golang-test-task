package storage

import (
	"app/internal/redis"
	"encoding/json"
	"fmt"
)

// SearchData searches data in Redis
func SearchData(globalID, mode, id string, redisClient redis.RedisClient) ([]map[string]interface{}, error) {
	var result []string
	var err error

	if globalID != "" {
		key := fmt.Sprintf("%s", globalID)
		result = append(result, key)
	} else {
		var tags []string

		if mode != "" {
			modeKey := fmt.Sprintf("index:mode:%s", mode)
			tags = append(tags, modeKey)
		}
		if id != "" {
			idKey := fmt.Sprintf("index:id:%s", id)
			tags = append(tags, idKey)
		}

		if len(tags) > 0 {
			result, err = redisClient.SInter(tags...)
			if err != nil {
				return nil, fmt.Errorf("failed to perform SINTER: %v", err)
			}
		}
	}

	var records []map[string]interface{}
	for _, recordKey := range result {
		data, err := redisClient.Get(recordKey)
		if err != nil {
			continue
		}

		var record map[string]interface{}
		err = json.Unmarshal([]byte(data), &record)
		if err != nil {
			continue
		}
		records = append(records, record)
	}

	return records, nil
}
