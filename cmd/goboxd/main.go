package main

import (
	"log"
	"net/http"

	"github.com/piyushiitm/goboxd/internal/api"
)

func main() {
	http.HandleFunc("/healthz", api.Healthz)
	http.HandleFunc("/run", api.Run)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
