package main

import (
	"deepwork/cmd"
)

func main() {
	cmd.Execute()
	// configPath, uid, gid := getRealUserHome()
	//
	// ensureConfigDir(configPath, uid, gid)
	//
	// minutes := flag.Int("min", 1, "Minutes")
	// blocklistDir := flag.String("dir", configPath, "Blocklist file path")
	// flag.Parse()
	//
	// domainsToBlock, err :=  internal.LoadDomainsFromDir(*blocklistDir)
	// if err != nil {
	// 	log.Fatalf("Fatal error during loading domains: %v", err)
	// }
	//
	// fmt.Printf("Loaded %d domains from blocklist.\n", len(domainsToBlock))
	//
	// fmt.Println("Starting timer for", *minutes, "minutes")
	//
	// iface, err := internal.GetDefaultInterface()
	// if err != nil {
	// 	log.Printf("[Warning] Could not find active network interface: %v\n", err)
	// } else {
	// 	if err := internal.SetCustomDNS(iface); err != nil {
	// 		log.Printf("[Warning] Failed to set Cloudflare DNS on %s: %v\n", iface, err)
	// 	} else {
	// 		fmt.Printf("Global DNS override engaged on interface: %s\n", iface)
	// 	}
	// }
	//
	// if err := internal.AddDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// internal.FlushDNS()
	//
	// internal.RunTUI(*minutes)
	//
	// fmt.Println("Restoring DNS service")
	// if err := internal.RevertDNS(iface); err != nil {
	// 	log.Printf("[Error] Failed to restore original DNS: %v\n", err)
	// } else {
	// 	fmt.Println("Global DNS restored.")
	// }
	//
	// fmt.Println("Restoring /etc/hosts...")
	// if err := internal.RemoveDomainsRestriction("/etc/hosts", domainsToBlock); err != nil {
	// 	log.Fatal(err)
	// }
	// internal.FlushDNS()
	//
}
