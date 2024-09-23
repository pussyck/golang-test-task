package config

import (
	"os"
)

type Config struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	DataSourceURL string
}

// LoadConfig load config
func LoadConfig() *Config {
	return &Config{
		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		DataSourceURL: getEnv("DATA_SOURCE_URL", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
