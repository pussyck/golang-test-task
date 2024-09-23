package main

import (
	"app/config"
	"app/internal/handler"
	"app/internal/metrics"
	"app/internal/redis"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	redisClient := redis.NewRedisClient(cfg.RedisHost+":"+cfg.RedisPort, cfg.RedisPassword)
	internal.InitMetrics()

	h := handler.NewHandler(redisClient)

	http.Handle("/load-data", internal.MetricsMiddleware(http.HandlerFunc(h.LoadDataHandler)))
	http.Handle("/search", internal.MetricsMiddleware(http.HandlerFunc(h.GetParkingDataHandler)))
	http.Handle("/metrics", internal.HandleMetrics())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
