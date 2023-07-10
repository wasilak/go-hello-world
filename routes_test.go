package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := mux.NewRouter()

	r.NewRoute().HandlerFunc(healthHandler)

	r.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}

	var response HealthResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	expectedResponse := HealthResponse{
		Status: "ok",
	}

	if response != expectedResponse {
		t.Errorf("Did not get expected response, got %+v", response)
	}
}

func TestRootHandler(t *testing.T) {
	// Create a request to the rootHandler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	recorder := httptest.NewRecorder()

	// Call the rootHandler function directly
	rootHandler(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Unexpected status code. Expected: %d, Got: %d", http.StatusOK, recorder.Code)
	}

	// Parse the response body
	var response APIResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}
