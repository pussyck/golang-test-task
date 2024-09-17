package main

import (
	"app/internal"
	"log"
	"net/http"
)

func main() {
	internal.InitRedis()

	http.HandleFunc("/load-data", internal.LoadDataHandler)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
