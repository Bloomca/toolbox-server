package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const (
	maxShareDataBytes          = 8 * 1024
	maxCreateShareRequestBytes = 64 * 1024
)

type createShareRequest struct {
	Data *string `json:"data"`
}

func getShare(w http.ResponseWriter, r *http.Request) {
	log.Printf("get share: id=%q", r.PathValue("id"))
	notFound(w, r)
}

func createShare(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxCreateShareRequestBytes)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var request createShareRequest
		if err := decoder.Decode(&request); err != nil {
			writeCreateShareRequestError(w, err)
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			writeCreateShareRequestError(w, err)
			return
		}

		if request.Data == nil {
			writeShareError(w, http.StatusBadRequest, "data must be a string")
			return
		}

		data := *request.Data
		if len(data) > maxShareDataBytes {
			writeShareError(w, http.StatusRequestEntityTooLarge, "data exceeds the 8 KiB limit")
			return
		}
		if !json.Valid([]byte(data)) {
			writeShareError(w, http.StatusBadRequest, "data must contain valid JSON")
			return
		}

		id, err := randomID()
		if err != nil {
			log.Printf("generate share ID: %v", err)
			writeShareError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if _, err := db.ExecContext(
			r.Context(),
			"INSERT INTO shared_spins (id, json) VALUES (?, ?)",
			id,
			data,
		); err != nil {
			log.Printf("save shared spin: %v", err)
			writeShareError(w, http.StatusInternalServerError, "internal server error")
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
}

func writeCreateShareRequestError(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeShareError(w, http.StatusRequestEntityTooLarge, "request body too large")
		return
	}

	writeShareError(w, http.StatusBadRequest, "invalid request body")
}

func writeShareError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{
		Status: status,
		Error:  message,
	}); err != nil {
		log.Printf("write share error response: %v", err)
	}
}

func randomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
