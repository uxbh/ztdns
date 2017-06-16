// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// listinterfacesCmd represents the listinterfaces command
var listinterfacesCmd = &cobra.Command{
	Use:   "listinterfaces",
	Short: "List network interfaces",
	Long: `List Interfaces (ztdns listinterfaces) lists the available network interfaces 
to start the server on.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get a list of interfaces from net
		ints, err := net.Interfaces()
		if err != nil {
			fmt.Printf("error getting interfaces: %s", err)
		}
		for i, n := range ints {
			fmt.Printf("%d: %v\n", i, n.Name)
			// Get a list of ip address on the interface
			addrs, _ := n.Addrs()
			for i, a := range addrs {
				ip, _, err := net.ParseCIDR(a.String())
				if err != nil {
					continue
				}
				fmt.Printf("\t%d: %s\n", i, ip)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(listinterfacesCmd)
}
