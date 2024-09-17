package internal

import (
	"context"
	"encoding/json"
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
