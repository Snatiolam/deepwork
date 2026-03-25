package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func loadDomainsFromDir(dirPath string) ([]string ,error) {
	var allDomains []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not read blocklist directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		fullPath := filepath.Join(dirPath, entry.Name())

		domains, err := parseBlocklistFile(fullPath)

		if err != nil {
			fmt.Printf("Warning: skipping file %s due to error: %v.\n", entry.Name(), err)
		}

		allDomains = append(allDomains, domains...)
	}

	return allDomains, nil
}

