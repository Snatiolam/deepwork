package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"deepwork/internal"

	"github.com/spf13/cobra"
)

var minutes int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a deepwork session",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		username := os.Getenv("SUDO_USER")
		if username == "" {
			return fmt.Errorf("this command must be run with sudo")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		username := os.Getenv("SUDO_USER")
		u, err := user.Lookup(username)
		if err != nil {
			return fmt.Errorf("could not find user info for %s: %w", username, err)
		}

		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)

		dir := blocklistDir
		if dir == "" {
			dir = filepath.Join(u.HomeDir, ".config", "deepwork", "blocklists")
		}

		internal.EnsureConfigDir(dir, uid, gid)

		domainsToBlock, err := internal.LoadDomainsFromDir(dir)
		if err != nil {
			return fmt.Errorf("fatal error during loading domains: %w", err)
		}

		fmt.Printf("Loaded %d domains from blocklist.\n", len(domainsToBlock))
		fmt.Println("Starting timer for", minutes, "minutes")

		iface, err := internal.GetDefaultInterface()
		if err != nil {
			log.Printf("[Warning] Could not find active network interface: %v\n", err)
		} else {
			if err := internal.SetCustomDNS(iface); err != nil {
				log.Printf("[Warning] Failed to set Cloudflare DNS on %s: %v\n", iface, err)
			} else {
				fmt.Printf("Global DNS override engaged on interface: %s\n", iface)
			}
		}

		if err := internal.AddDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
			return err
		}

		internal.FlushDNS()

		if err := internal.RunTUI(minutes); err != nil {
			return err
		}

		fmt.Println("Restoring DNS service")
		if err := internal.RevertDNS(iface); err != nil {
			log.Printf("[Error] Failed to restore original DNS: %v\n", err)
		} else {
			fmt.Println("Global DNS restored.")
		}

		fmt.Println("Restoring /etc/hosts...")
		if err := internal.RemoveDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
			return err
		}
		internal.FlushDNS()

		return nil
	},
}

func init() {
	startCmd.Flags().IntVarP(&minutes, "min", "m", 1, "Minutes to run")
	rootCmd.AddCommand(startCmd)
}
