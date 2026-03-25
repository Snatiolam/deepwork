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

	if err := addDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
		log.Fatal(err)
	}

	flusDNS()

	fmt.Println("The program is running. Press Ctrl-C to exit.")

	runTUI(*minutes)

	fmt.Println("Restoring /etc/hosts...")
	if err := removeDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
		log.Fatal(err)
	}

	flusDNS()
}
