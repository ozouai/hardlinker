package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// RequestBody represents the expected JSON request body
type RequestBody struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	// Parse command line arguments
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <source> <destination>", os.Args[0])
	}

	source := os.Args[1]
	destination := os.Args[2]

	// Create request body
	requestBody := RequestBody{
		Source:      source,
		Destination: destination,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Determine the server address
	addr := os.Getenv("LINKER_ADDR")
	if addr == "" {
		addr = "localhost:5070"
	}

	// Make HTTP POST request
	resp, err := http.Post(
		fmt.Sprintf("http://%s/link", addr),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		// For non-200 responses, print the error message from the server
		log.Fatalf("Server error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// If we get here, the request was successful
	fmt.Printf("Successfully created hard link from %s to %s\n", source, destination)
}
