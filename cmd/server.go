// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uxbh/ztdns/dnssrv"
	"github.com/uxbh/ztdns/ztapi"
	"github.com/aleacevedo/ztdns/helpers"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run ztDNS server",
	Long: `Server (ztdns server) will start the DNS server.append
	
	Example: ztdns server`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check config and bail if anything important is missing.
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Setting Debug Mode")
		}
		if viper.GetString("ZT.API") == "" {
			return fmt.Errorf("no API key provided")
		}
		if len(viper.GetStringMapString("Networks")) == 0 {
			return fmt.Errorf("no Domain / Network ID pairs Provided")
		}
		if viper.GetString("ZT.URL") == "" {
			return fmt.Errorf("no URL provided. Run ztdns mkconfig first")
		}
		if viper.GetString("suffix") == "" {
			return fmt.Errorf("no DNS Suffix provided. Run ztdns mkconfig first")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Update the DNSDatabase
		lastUpdate := updateDNS()
		req := make(chan string)
		// Start the DNS server
		go dnssrv.Start(viper.GetString("interface"), viper.GetInt("port"), viper.GetString("suffix"), req)

		refresh := viper.GetInt("DbRefresh")
		if refresh == 0 {
			refresh = 30
		}
		for {
			// Block until a new request comes in
			n := <-req
			log.Debugf("Got request for %s", n)
			// If the database hasn't been updated in the last "refresh" minutes, update it.
			if time.Since(lastUpdate) > time.Duration(refresh)*time.Minute {
				log.Infof("DNSDatabase is stale. Refreshing.")
				lastUpdate = updateDNS()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().String("interface", "", "interface to listen on")
	viper.BindPFlag("interface", serverCmd.PersistentFlags().Lookup("interface"))
}

func updateDNS() time.Time {
	// Get config info
	API := viper.GetString("ZT.API")
	URL := viper.GetString("ZT.URL")
	suffix := viper.GetString("suffix")

	// Get all configured networks:
	for domain, id := range viper.GetStringMapString("Networks") {
		// Get ZeroTier Network info
		ztnetwork, err := ztapi.GetNetworkInfo(API, URL, id)
		if err != nil {
			log.Fatalf("Unable to update DNS entries: %s", err.Error())
		}

		// Get list of members in network
		log.Infof("Getting Members of Network: %s (%s)", ztnetwork.Config.Name, domain)
		members, err := ztapi.GetMemberList(API, URL, ztnetwork.ID)
		if err != nil {
			log.Fatalf("Unable to update DNS entries: %s", err.Error())
		}
		log.Infof("Got %d members", len(*members))

		for _, member := range *members {
			// For all online members
			if member.Online {

				subdomains := viper.GetStringSlice(fmt.Sprintf("Subdomains.%s.%s", domain, member.Name))
				
				routes := helpers.GenerateRoutes(member.Name, domain, suffix, subdomains)
				ip6, ip4 := helpers.GenerateIPs(ztnetwork, &member)

				for _, route := range routes {
					// Clear current DNS records
					dnssrv.DNSDatabase[route] = dnssrv.Records{}
					// Add the record to the database
					log.Infof("Updating %-15s IPv4: %-15s IPv6: %s", route, ip4, ip6)
					dnssrv.DNSDatabase[route] = dnssrv.Records{
						A:    ip4,
						AAAA: ip6,
					}
				}
			}
		}
	}

	// Return the current update time
	return time.Now()
}
