// Package cmd contains the CLI command logic
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	blocklistDir string
)

var rootCmd = &cobra.Command{
	Use: "focus",
	Short: "A high-intensity deepwork tool for linux",
	Long: "A tool to help you maintain deep focus by blocking distracting websites.",
}

func init() {
	defaultPath := getDefaultBlocklistPath()
	rootCmd.PersistentFlags().StringVarP(&blocklistDir, "dir", "d", defaultPath, "Blocklist directory path")
}

func getDefaultBlocklistPath() string {
	username := os.Getenv("SUDO_USER")
	if username == "" {
		return ""
	}
	u, err := user.Lookup(username)
	if err != nil {
		return ""
	}
	return filepath.Join(u.HomeDir, ".config", "deepwork", "blocklists")
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
