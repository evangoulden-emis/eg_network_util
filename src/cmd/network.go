package main

import (
	"fmt"
	"net"
	"strings"
)

// getPhysicalInterface returns the IP of the first physical network interface
func getPhysicalInterface() (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    
    for _, iface := range interfaces {
        // Skip virtual interfaces
        if strings.Contains(iface.Name, "lo") || // loopback
           strings.Contains(iface.Name, "tun") || // VPN tunnel
           strings.Contains(iface.Name, "tap") || // VPN tap
           strings.Contains(iface.Name, "ppp") || // PPP
           strings.Contains(iface.Name, "utun") || // macOS VPN
           strings.Contains(iface.Flags.String(), "loopback") ||
           (iface.Flags & net.FlagUp) == 0 { // interface not up
            continue
        }
        
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }
        
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok {
                if ip := ipnet.IP.To4(); ip != nil {
                    return ip.String(), nil
                }
            }
        }
    }
    return "", fmt.Errorf("no suitable physical interface found")
}