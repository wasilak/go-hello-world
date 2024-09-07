package gorilla

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	// Create a mock handler to be wrapped by the middlewares
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something in the handler
		w.WriteHeader(http.StatusOK)
	})

	// Define some mock middlewares
	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Do middleware 1 things
			// ...

			// Call the next middleware/handler in chain
			next(w, r)
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Do middleware 2 things
			// ...

			// Call the next middleware/handler in chain
			next(w, r)
		}
	}

	// Apply the middlewares using the Chain function
	handler := chain(mockHandler, middleware1, middleware2)

	// Create a test server with the wrapped handler
	server := httptest.NewServer(handler)
	defer server.Close()

	// Send a request to the server
	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code. Expected: %d, Got: %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a mock handler to be wrapped by the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something in the handler
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with the middleware applied
	handler := logging()(mockHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Send a request to the server
	client := &http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code. Expected: %d, Got: %d", http.StatusOK, resp.StatusCode)
	}
}
