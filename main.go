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

	if err := http.ListenAndServe(address, server()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func server() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("POST /api/spinny/share", createShare)
	mux.HandleFunc("GET /api/spinny/share/{id}", getShare)
	mux.HandleFunc("/", notFound)

	return mux
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("{\"status\":404,\"error\":\"not found\"}\n"))
}
