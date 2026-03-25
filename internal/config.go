package internal

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func getRealUserHome() (string, int, int) {
	username := os.Getenv("SUDO_USER")
	if username == "" {
		log.Fatal("This program must be run with sudo")
	}

	u, err := user.Lookup(username)
	if err != nil {
		log.Fatalf("Could not find user info for %s: %v", username, err)
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	path := filepath.Join(u.HomeDir, ".config", "deepwork", "blocklists")

	return path, uid, gid
}

func EnsureConfigDir(path string, uid, gid int) {
	if _, err := os.Stat(path); err == nil {
		return
	}

	// fmt.Printf("Creating config directory: %s\n", path)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	parent := filepath.Dir(path)
	os.Chown(parent, uid, gid)
	os.Chown(path, uid, gid)
}
