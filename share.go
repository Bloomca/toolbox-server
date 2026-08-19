package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

func getShare(w http.ResponseWriter, r *http.Request) {
	log.Printf("get share: id=%q", r.PathValue("id"))
	notFound(w, r)
}

func createShare(w http.ResponseWriter, _ *http.Request) {
	id, err := randomID()
	if err != nil {
		log.Printf("generate share ID: %v", err)
		http.Error(w, "failed to generate id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Status int    `json:"status"`
		ID     string `json:"id"`
	}{
		Status: http.StatusOK,
		ID:     id,
	}); err != nil {
		log.Printf("write share response: %v", err)
	}
}

func randomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
