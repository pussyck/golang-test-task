package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis(addr, password string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

func SAdd(key string, value string) error {
	return redisClient.SAdd(ctx, key, value).Err()
}

func SInter(keys ...string) ([]string, error) {
	return redisClient.SInter(ctx, keys...).Result()
}

func SMembers(key string) ([]string, error) {
	return redisClient.SMembers(ctx, key).Result()
}

func Set(key string, value interface{}) error {
	return redisClient.Set(ctx, key, value, 0).Err()
}

func HGetAll(key string) (map[string]string, error) {
	return redisClient.HGetAll(ctx, key).Result()
}

func Get(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}
