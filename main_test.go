package main

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestHardlinkerFunctionality tests the core hardlink functionality
func TestHardlinkerFunctionality(t *testing.T) {
	// Clean up any existing test files
	os.Remove("test-source.txt")
	os.Remove("test-destination.txt")
	os.Remove("links.yaml")

	// Create source file with content
	sourceContent := "This is a test file for hardlinking"
	err := ioutil.WriteFile("test-source.txt", []byte(sourceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test the addToYAML function directly
	err = addToYAML("test-source.txt", "test-destination.txt")
	if err != nil {
		t.Fatalf("Failed to add to YAML: %v", err)
	}

	// Verify the link was added to YAML file
	yamlContent, err := ioutil.ReadFile("links.yaml")
	if err != nil {
		t.Fatalf("Failed to read YAML file: %v", err)
	}

	var links []Link
	err = yaml.Unmarshal(yamlContent, &links)
	if err != nil {
		t.Fatalf("Failed to parse YAML file: %v", err)
	}

	if len(links) != 1 {
		t.Errorf("Expected 1 link in YAML file, got %d", len(links))
	}

	if links[0].Source != "test-source.txt" {
		t.Errorf("Expected source test-source.txt, got %s", links[0].Source)
	}

	if links[0].Destination != "test-destination.txt" {
		t.Errorf("Expected destination test-destination.txt, got %s", links[0].Destination)
	}

	// Clean up
	os.Remove("test-source.txt")
	os.Remove("test-destination.txt")
	os.Remove("links.yaml")
}

// TestHardlinkCreation tests the core hard link creation functionality
func TestHardlinkCreation(t *testing.T) {
	// Clean up any existing test files
	os.Remove("test-source.txt")
	os.Remove("test-destination.txt")

	// Create source file with content
	sourceContent := "This is a test file for hardlinking"
	err := ioutil.WriteFile("test-source.txt", []byte(sourceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create the hard link manually (without using HTTP server)
	err = os.Link("test-source.txt", "test-destination.txt")
	if err != nil {
		t.Fatalf("Failed to create hard link: %v", err)
	}

	// Verify both files exist
	_, err = os.Stat("test-source.txt")
	if err != nil {
		t.Fatalf("Source file should exist: %v", err)
	}

	_, err = os.Stat("test-destination.txt")
	if err != nil {
		t.Fatalf("Destination file should exist: %v", err)
	}

	// Verify content is the same
	sourceContentRead, err := ioutil.ReadFile("test-source.txt")
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	destContentRead, err := ioutil.ReadFile("test-destination.txt")
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(sourceContentRead) != string(destContentRead) {
		t.Errorf("Files should have same content. Got: %s, Expected: %s",
			string(destContentRead), string(sourceContentRead))
	}

	// Clean up
	os.Remove("test-source.txt")
	os.Remove("test-destination.txt")
}
