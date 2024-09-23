package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient interface {
	Set(key string, value interface{}) error
	SAdd(key string, value string) error
	SInter(keys ...string) ([]string, error)
	SMembers(key string) ([]string, error)
	HGetAll(key string) (map[string]string, error)
	Get(key string) (string, error)
}

var ctx = context.Background()

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &Client{client: client}
}

func (r *Client) Set(key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *Client) SAdd(key string, value string) error {
	return r.client.SAdd(ctx, key, value).Err()
}

func (r *Client) SInter(keys ...string) ([]string, error) {
	return r.client.SInter(ctx, keys...).Result()
}

func (r *Client) SMembers(key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *Client) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *Client) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
