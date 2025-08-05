package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	yaml "gopkg.in/yaml.v3"
)

// Global mutex for protecting YAML file access
var yamlMutex = sync.Mutex{}

// Link represents a source and destination for a hard link
type Link struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// RequestBody represents the expected JSON request body
type RequestBody struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	http.HandleFunc("/link", linkHandler)

	log.Println("Starting hardlinker server on port 5070")
	err := http.ListenAndServe(":5070", nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var reqBody RequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reqBody.Source == "" || reqBody.Destination == "" {
		http.Error(w, "Source and destination are required", http.StatusBadRequest)
		return
	}

	// Check if destination already exists
	if _, err := os.Stat(reqBody.Destination); err == nil {
		http.Error(w, "Destination already exists", http.StatusConflict)
		return
	}

	// Create destination directory if needed
	destDir := filepath.Dir(reqBody.Destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		http.Error(w, "Failed to create destination directory", http.StatusInternalServerError)
		return
	}

	// Create hard link
	err = os.Link(reqBody.Source, reqBody.Destination)
	if err != nil {
		// Check if the error is because source doesn't exist
		if _, statErr := os.Stat(reqBody.Source); statErr != nil {
			http.Error(w, "Source file does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create hard link: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add link to YAML file
	if err := addToYAML(reqBody.Source, reqBody.Destination); err != nil {
		log.Printf("Warning: Failed to add link to YAML file: %v", err)
		// Note: We don't return an error for YAML failure as the hard link was successful
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "success", "source": "%s", "destination": "%s"}`, reqBody.Source, reqBody.Destination)
}

func addToYAML(source, destination string) error {
	// Determine YAML file path
	yamlPath := "links.yaml"
	if envPath := os.Getenv("LINK_YAML"); envPath != "" {
		yamlPath = envPath
	}

	// Lock the YAML file for concurrent access protection
	yamlMutex.Lock()
	defer yamlMutex.Unlock()

	// Read existing links from YAML file
	var links []Link
	if _, err := os.Stat(yamlPath); err == nil {
		// File exists, read it
		data, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return fmt.Errorf("failed to read YAML file: %w", err)
		}

		if err := yaml.Unmarshal(data, &links); err != nil {
			return fmt.Errorf("failed to parse YAML file: %w", err)
		}
	}

	// Add new link
	newLink := Link{
		Source:      source,
		Destination: destination,
	}
	links = append(links, newLink)

	// Write back to YAML file
	data, err := yaml.Marshal(&links)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := ioutil.WriteFile(yamlPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}
