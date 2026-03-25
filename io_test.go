package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadDomains(t *testing.T) {
	tmpDir := t.TempDir()

	mockContent := "youtube.com\n# This is a comment\n\ntwitter.com\n   reddit.com   \n"
	mockFile := filepath.Join(tmpDir, "social.txt")

	err := os.WriteFile(mockFile, []byte(mockContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock blocklist: %v", err)
	}

	domains, err := loadDomainsFromDir(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := []string{"youtube.com", "twitter.com", "reddit.com"}

	if !reflect.DeepEqual(domains, expected) {
		t.Errorf("Expected %v, got %v", domains, expected)
	}
}
