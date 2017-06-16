// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mkconfigCmd represents the mkconfig command
var mkconfigCmd = &cobra.Command{
	Use:   "mkconfig",
	Short: "Make a new config file",
	Long: `mkconfig (ztdns mkconfig) creates a new configuation file.
If you do not specify a filename the default is ./.ztdns.toml

Example: ztdns mkconfig .filenames.toml`,
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

# This section contains information related to your ZeroTier config
[ZT]
# API is used to contact the ZeroTier controller API service.
API = ""
# Network is your ZeroTier network ID.
Network = ""
# URL is the url of the ZeroTier controller API
URL = "https://my.zerotier.com/api"
`)
		}
	},
}

func init() {
	RootCmd.AddCommand(mkconfigCmd)
}
