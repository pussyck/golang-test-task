package main

import (
	"app/internal"
	"log"
	"net/http"
)

func main() {
	internal.InitRedis()

	internal.InitMetrics()

	http.Handle("/load-data", internal.MetricsMiddleware(http.HandlerFunc(internal.LoadDataHandler)))
	http.Handle("/search", internal.MetricsMiddleware(http.HandlerFunc(internal.GetParkingDataHandler)))
	http.Handle("/metrics", internal.HandleMetrics())

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
