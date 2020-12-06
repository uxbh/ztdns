// Copyright © 2017 aleacevedo
// This file is part of github.com/aleacevedo/ztdns.

// Package helpers contain some helpers to improve route generation
package helpers

import (
	"fmt"
	"net"

	"github.com/uxbh/ztdns/ztapi"
)

func GetRoot(name string, network string, suffix string) string {
	if (network != "default") {
		return  fmt.Sprintf("%s.%s.%s.", name, network, suffix)
	}
	return fmt.Sprintf("%s.%s.", name, suffix)
}

func GenerateRoutes(name string, network string, suffix string, subdomains []string) []string {
	routes := []string{}
	root := GetRoot(name, network, suffix)
	for _, subdomain := range subdomains {
		routes = append(routes, fmt.Sprintf("%s.%s", subdomain, root))
	}
	routes = append(routes, root)
	return routes
}

func GenerateIPs(ztnetwork *ztapi.Network, member *ztapi.Member) ([]net.IP, []net.IP) {
	ip6 := []net.IP{}
	ip4 := []net.IP{}
	// Get 6Plane address if network has it enabled
	if ztnetwork.Config.V6AssignMode.Sixplane {
		ip6 = append(ip6, member.Get6Plane())
	}
	// Get RFC4193 address if network has it enabled
	if ztnetwork.Config.V6AssignMode.Rfc4193 {
		ip6 = append(ip6, member.GetRFC4193())
	}

	// Get the rest of the address assigned to the member
	for _, a := range member.Config.IPAssignments {
		ip4 = append(ip4, net.ParseIP(a))
	}

	return ip6, ip4
}