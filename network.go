package main

import (
	"fmt"
	"net"
	"os/exec"
)

func getDefaultInterface() (string, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "", fmt.Errorf("could not determine default route: %w", err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.Equal(localAddr.IP) {
				return iface.Name, nil
			}
		}
	}
	return "", fmt.Errorf("could not find interface for IP %s", localAddr.IP)
}

func setCustomDNS(iface string) error {
	cmd := exec.Command("resolvectl", "dns", iface, "1.1.1.3", "1.0.0.3")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set Cloudflare DNS via resolvectl: %w", err)
	}
	return nil
}

func revertDNS(iface string) error {
	cmd := exec.Command("resolvectl", "revert", iface)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to revert DNS via resolvectl: %w", err)
	}
	return nil
}
