package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Setup test environment here if needed
	flag.Parse()
	exitCode := m.Run()
	// Teardown test environment here if needed
	os.Exit(exitCode)
}

func TestMainFunction(t *testing.T) {
	// Save and restore original command-line arguments
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Save and restore original flag values
	originalListenAddr := listenAddr
	originalLogLevel := logLevel
	originalLogFormat := logFormat
	originalOtelEnabled := otelEnabled

	defer func() {
		listenAddr = originalListenAddr
		logLevel = originalLogLevel
		logFormat = originalLogFormat
		otelEnabled = originalOtelEnabled
	}()

	// Set desired flag values for testing
	listenAddr = ":0" // Use an available port for testing
	logLevel = "info"
	logFormat = "text"
	otelEnabled = true

	// Start the server in a separate goroutine
	go func() {
		main() // Call the main function
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Make a request to the server
	resp, err := http.Get(fmt.Sprintf("http://%s", listenAddr))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
}
