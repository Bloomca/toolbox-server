package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	address := ":" + port
	log.Printf("listening on %s", address)

	if err := http.ListenAndServe(address, http.NewServeMux()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
