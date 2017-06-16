// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

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
		ints, _ := net.Interfaces()

		for i, n := range ints {
			addrs, _ := n.Addrs()
			fmt.Printf("%d: %v\n", i, n.Name)
			for i, a := range addrs {
				ip, _, err := net.ParseCIDR(a.String())
				if err != nil {
					return
				}
				fmt.Printf("\t%d: %s\n", i, ip)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(listinterfacesCmd)
}
