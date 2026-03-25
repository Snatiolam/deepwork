package main

import (
	"os"
	"strings"
	"testing"
)

func TestHostsManipulation(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "fake_hosts")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tmpFile.Name())

	initialContent := "127.0.0.1   localhost\n::1    localhost\n"
	os.WriteFile(tmpFile.Name(), []byte(initialContent), 0644)

	domainsToBlock := []string{"distract.com", "news.com"}

	err = addDomainsRestriction(tmpFile.Name(), domainsToBlock)
	if err != nil {
		t.Fatalf("Failed to add restrictions: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpFile.Name())
	content := string(contentBytes)

	if !strings.Contains(content, "0.0.0.0 distract.com") {
		t.Errorf("Fake hosts file is missing the IPv4 block for distraction.com")
	}
	if !strings.Contains(content, ":: news.com") {
		t.Errorf("Fake hosts file is missing the IPv6 block for news.com")
	}

	err = removeDomainsRestriction(tmpFile.Name(), domainsToBlock)
	if err != nil {
		t.Fatalf("Failed to remove restrictions: %v", err)
	}

	finalContentBytes, _ := os.ReadFile(tmpFile.Name())
	finalContent := string(finalContentBytes)

	if strings.Contains(finalContent, "distraction.com") {
		t.Errorf("distraction.com was not properly removed from hosts")
	}

	if !strings.Contains(finalContent, "127.0.0.1   localhost") {
		t.Errorf("FATAL: The removal function destroyed the original system hosts data")
	}
}

func TestAddDomainsRestriction(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "fake_hosts")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		domains []string
	}{
		{"Single domain", []string{"reddit.com"}},
		{"Multiple domains", []string{"twitter.com", "youtube.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addDomainsRestriction(tmpFile.Name(), tt.domains)
			if err != nil {
				t.Errorf("failed to add restriction: %v", err)
			}

			content, _ := os.ReadFile(tmpFile.Name())
			for _, d := range tt.domains {
				if !contains(string(content), d) {
					t.Errorf("expected domain %s not found in file", d)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
