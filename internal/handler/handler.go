package handler

import "app/internal/redis"

type Handler struct {
	RedisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{RedisClient: redisClient}
}
