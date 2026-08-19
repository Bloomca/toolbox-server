package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestShareEndToEnd(t *testing.T) {
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	testServer := httptest.NewServer(server(db))
	defer testServer.Close()

	data := `{"id":"list-1","title":"Test list","choices":[{"id":"sun","label":"Sun","weight":1,"included":true,"parentChoiceId":null}]}`
	requestBody, err := json.Marshal(struct {
		Data string `json:"data"`
	}{Data: data})
	if err != nil {
		t.Fatalf("encode create request: %v", err)
	}

	createResponse, err := testServer.Client().Post(
		testServer.URL+"/api/spinny/share",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		t.Fatalf("create shared spin: %v", err)
	}
	defer createResponse.Body.Close()

	if createResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(createResponse.Body)
		t.Fatalf("expected create status %d, got %d: %s", http.StatusOK, createResponse.StatusCode, body)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(createResponse.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected create response to contain an ID")
	}

	getResponse, err := testServer.Client().Get(testServer.URL + "/api/spinny/share/" + created.ID)
	if err != nil {
		t.Fatalf("get shared spin: %v", err)
	}
	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(getResponse.Body)
		t.Fatalf("expected get status %d, got %d: %s", http.StatusOK, getResponse.StatusCode, body)
	}

	var fetched struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}
	if err := json.NewDecoder(getResponse.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if fetched.ID != created.ID {
		t.Fatalf("expected ID %q, got %q", created.ID, fetched.ID)
	}
	if fetched.Data != data {
		t.Fatalf("expected data %q, got %q", data, fetched.Data)
	}
}
