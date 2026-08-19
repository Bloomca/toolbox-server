package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()

	server().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if response.Body.String() != "ok\n" {
		t.Fatalf("expected body %q, got %q", "ok\n", response.Body.String())
	}
}

func TestNotFound(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/does-not-exist", nil)
	response := httptest.NewRecorder()

	server().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}
