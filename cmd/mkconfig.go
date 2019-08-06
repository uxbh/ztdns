// Copyright © 2017 uxbh
// This file is part of github.com/hatemosphere/ztdns.

package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mkconfigCmd represents the mkconfig command
var mkconfigCmd = &cobra.Command{
	Use:   "mkconfig",
	Short: "Make a new config file",
	Long: `mkconfig (ztdns mkconfig) creates a new configuation file.
If you do not specify a filename the default is ./.ztdns.toml

Example: ztdns mkconfig [.filename.toml]`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := "./.ztdns.toml"
		if len(args) > 0 {
			filename = args[0]
		}
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("Creating new config file in %s", filename)
			file, err := os.Create(filename)
			if err != nil {
				log.Fatalf("Could not create file: %s", err.Error())
			}
			defer file.Close()
			file.WriteString(`# Configuration file for ztDNS

suffix = "zt"
port = 53
interface = "zt0"

# Number of minutes to wait before updating the DNS database again (Default: 30)
DBRefresh = 30

# This section contains information related to your ZeroTier config
[ZT]
# API is used to contact the ZeroTier controller API service.
API = ""
# URL is the url of the ZeroTier controller API
URL = "https://my.zerotier.com/api"

# This section contains one or more ZeroTier networks
# Format is: domain = "NetworkID"
# Domain does not have to match the configured network name
[Networks]

# Match section contains zero or more match pairs to create Round robin dns
# Format is: "regexp to match hosts" = "hostname"
# Example 1:
#   "k8s-node-\w" = "k8s-nodes"
#   From nodes with names k8s-node-23refw, k8s-node-09sf8g
#   will create round robin record k8s-nodes
[RoundRobin]

`)
		}
	},
}

func init() {
	RootCmd.AddCommand(mkconfigCmd)
}
