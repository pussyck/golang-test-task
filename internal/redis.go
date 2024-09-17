package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

func processJSON(data []byte) error {
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	for _, record := range records {
		globalID, ok := record["global_id"].(float64)
		if !ok {
			return fmt.Errorf("invalid global_id format")
		}
		globalIDStr := fmt.Sprintf("%.0f", globalID)

		recordJSON, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %v", err)
		}
		if err := redisClient.Set(ctx, globalIDStr, recordJSON, 0).Err(); err != nil {
			return fmt.Errorf("failed to save data to Redis: %v", err)
		}
	}

	return nil
}

func GetParkingDataByField(field, value string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if field == "global_id" {
		val, err := redisClient.Get(ctx, value).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("no data found for %s", value)
			}
			return nil, fmt.Errorf("failed to get data from Redis: %v", err)
		}

		var record map[string]interface{}
		if err := json.Unmarshal([]byte(val), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal record: %v", err)
		}
		result = append(result, record)
	} else {
		keys, err := redisClient.Keys(ctx, "*").Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			jsonData, err := redisClient.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}

			var record map[string]interface{}
			if err := json.Unmarshal([]byte(jsonData), &record); err != nil {
				return nil, err
			}

			if CompareValue(record, field, value) {
				result = append(result, record)
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no data found for %s", value)
	}

	return result, nil
}

func CompareValue(record map[string]interface{}, field string, value string) bool {
	val, ok := record[field]
	if !ok {
		return false
	}

	valStr := fmt.Sprintf("%v", val)
	valueStr := fmt.Sprintf("%v", value)

	return valStr == valueStr
}
