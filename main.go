package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	configPath, uid, gid := getRealUserHome()

	ensureConfigDir(configPath, uid, gid)

	minutes := flag.Int("min", 1, "Minutes")
	blocklistDir := flag.String("dir", configPath, "Blocklist file path")
	flag.Parse()

	domainsToBlock, err := loadDomainsFromDir(*blocklistDir)
	if err != nil {
		log.Fatalf("Fatal error during loading domains: %v", err)
	}

	fmt.Printf("Loaded %d domains from blocklist.\n", len(domainsToBlock))

	fmt.Println("Starting timer for", *minutes, "minutes")

	iface, err := getDefaultInterface()
	if err != nil {
		log.Printf("[Warning] Could not find active network interface: %v\n", err)
	} else {
		if err := setCustomDNS(iface); err != nil {
			log.Printf("[Warning] Failed to set Cloudflare DNS on %s: %v\n", iface, err)
		} else {
			fmt.Printf("Global DNS override engaged on interface: %s\n", iface)
		}
	}

	if err := addDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
		log.Fatal(err)
	}

	flushDNS()

	runTUI(*minutes)

	fmt.Println("Restoring DNS service")
	if err := revertDNS(iface); err != nil {
		log.Printf("[Error] Failed to restore original DNS: %v\n", err)
	} else {
		fmt.Println("Global DNS restored.")
	}

	fmt.Println("Restoring /etc/hosts...")
	if err := removeDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
		log.Fatal(err)
	}
	flushDNS()
}
