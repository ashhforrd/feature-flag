package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ashhforrd/feature-flags/services/config-service/internal/flags"
)

func main() {
	repo := flags.NewMemoryRepository()
	flagHandler := flags.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	flagHandler.RegisterRoutes(mux)

	log.Println(("config service listening on :8080"))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
