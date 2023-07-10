package main

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	handler := Chain(mockHandler, middleware1, middleware2)

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
	handler := Logging()(mockHandler)
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

func TestGenerateKey(t *testing.T) {
	t.Run("Valid Key Generation", func(t *testing.T) {
		key, err := GenerateKey()
		assert.NoError(t, err, "Expected no error during key generation")
		assert.Len(t, key, 64, "Expected key length of 64 characters")
	})

	t.Run("Randomness", func(t *testing.T) {
		// Generate multiple keys and verify they are different
		keys := make(map[string]struct{})
		for i := 0; i < 100; i++ {
			key, err := GenerateKey()
			assert.NoError(t, err, "Expected no error during key generation")
			_, exists := keys[key]
			assert.False(t, exists, "Expected unique keys")
			keys[key] = struct{}{}
		}
	})

	t.Run("Key Length", func(t *testing.T) {
		key, err := GenerateKey()
		assert.NoError(t, err, "Expected no error during key generation")

		decoded, err := hex.DecodeString(key)
		assert.NoError(t, err, "Expected no error during key decoding")
		assert.Len(t, decoded, 32, "Expected key length of 32 bytes")
	})
}
