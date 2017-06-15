// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/uxbh/ztdns/dnssrv"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run ztDNS server",
	Long: `Server (ztdns server) will start the DNS server.append
	
	Example: ztdns server`,
	Run: func(cmd *cobra.Command, args []string) {
		dnssrv.Start(viper.GetInt("port"), viper.GetString("suffix"))
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
