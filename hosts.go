package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func parseBlocklistFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open blocklist: %w", err)
	}

	defer file.Close()

	var domains []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		domains = append(domains, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading blocklist: %w", err)
	}

	return domains, nil
}

func addDomainsRestriction(path string, domains []string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	defer file.Close()

	file.WriteString("\n")

	for _, domain := range domains {
		if _, err = file.WriteString("0.0.0.0 " + domain + "\n"); err != nil {
			return err
		}
		if _, err = file.WriteString(":: " + domain + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func removeDomainsRestriction(path string, domains []string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	for _, line := range lines {
		shouldKeep := true
		for _, domain := range domains {
			if strings.Contains(line, domain) {
				shouldKeep = false
				break
			}
		}

		if shouldKeep {
			newLines = append(newLines, line)
		}
	}

	if newLines[len(newLines)-1] == "" {
		newLines = newLines[:len(newLines)-1]
	}

	newContent := strings.Join(newLines, "\n")

	err = os.WriteFile(path, []byte(newContent), 0644)

	if err != nil {
		return err
	}

	return nil
}

func flushDNS() error {
	cmdArgs, err := getDNSFlushCommand()
	if err != nil {
		return fmt.Errorf("environment check failed: %w", err)
	}

	// fmt.Printf("Detected DNS manager: %s\n", cmdArgs[0])

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	if err := cmd.Run(); err != nil {
		return err
	}

	// fmt.Println("System DNS cache flushed successfully")
	return nil
}

func getDNSFlushCommand() ([]string, error) {
	supportedManagers := map[string][]string{
		"resolvectl":      {"flush-caches"},
		"systemd-resolve": {"--flush-caches"},
		"nscd":            {"-i", "hosts"},
	}

	for binary, args := range supportedManagers {
		path, err := exec.LookPath(binary)
		if err == nil {
			return append([]string{path}, args...), nil
		}
	}

	return nil, fmt.Errorf("no supported DNS manager found (tried: resolvectl, nscd, etc.)")
}
