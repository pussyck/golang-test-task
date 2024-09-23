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
	redis.InitRedis(cfg.RedisHost+":"+cfg.RedisPort, cfg.RedisPassword)

	internal.InitMetrics()

	http.Handle("/load-data", internal.MetricsMiddleware(http.HandlerFunc(handler.LoadDataHandler)))
	http.Handle("/search", internal.MetricsMiddleware(http.HandlerFunc(handler.GetParkingDataHandler)))
	http.Handle("/metrics", internal.HandleMetrics())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
