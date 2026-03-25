package main

import (
	"os"
	"strings"
	"testing"
)

func TestAddDomainsRestriction(t *testing.T) {
	// 1. Setup a temporary fake hosts file
	tmpFile, err := os.CreateTemp("", "fake_hosts")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// 2. Define test cases
	tests := []struct {
		name    string
		domains []string
	}{
		{"Single domain", []string{"reddit.com"}},
		{"Multiple domains", []string{"twitter.com", "youtube.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function pointing to our TEMP file
			err := addDomainsRestriction(tmpFile.Name(), tt.domains)
			if err != nil {
				t.Errorf("failed to add restriction: %v", err)
			}

			// Read it back and verify
			content, _ := os.ReadFile(tmpFile.Name())
			for _, d := range tt.domains {
				if !contains(string(content), d) {
					t.Errorf("expected domain %s not found in file", d)
				}
			}
		})
	}
}

// Helper for testing
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
