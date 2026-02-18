package main

import (
	"fmt"
	"log"
	"net/http"

	"rizon-test-task/internal/config"
)

func main() {
	cfg := config.GetServerConfig()

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	fmt.Println("Endpoints:")
	fmt.Printf("  GET http://localhost:%s/health\n", cfg.Port)

	if err := http.ListenAndServe(cfg.Addr(), nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
