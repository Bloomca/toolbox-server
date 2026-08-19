package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		log.Fatal("DATABASE_PATH environment variable is required")
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	address := ":" + port
	log.Printf("listening on %s", address)

	if err := http.ListenAndServe(address, server(db)); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func server(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("POST /api/spinny/share", createShare(db))
	mux.HandleFunc("GET /api/spinny/share/{id}", getShare(db))
	mux.HandleFunc("/", notFound)

	return mux
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("{\"status\":404,\"error\":\"not found\"}\n"))
}
