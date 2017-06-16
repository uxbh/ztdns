// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package cmd

import (
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/uxbh/ztdns/dnssrv"
	"gitlab.com/uxbh/ztdns/ztapi"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run ztDNS server",
	Long: `Server (ztdns server) will start the DNS server.append
	
	Example: ztdns server`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Setting Debug Mode")
		}

		if viper.GetString("ZT.API") == "" {
			log.Fatal("No API key provided")
		}
		if viper.GetString("ZT.URL") == "" {
			log.Fatal("No URL provided. Run ztdns mkconfig first")
		}
		log.Debugf("Using API: %s", viper.GetString("ZT.API"))
		if viper.GetString("ZT.Network") == "" {
			log.Fatal("No Network ID Provided")
		}
		lastUpdate := updateDNS()
		req := make(chan bool)
		go dnssrv.Start(viper.GetString("interface"), viper.GetInt("port"), viper.GetString("suffix"), req)
		for {
			<-req
			log.Debug("Got Request")
			if time.Since(lastUpdate) > 30*time.Minute {
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
	API := viper.GetString("ZT.API")
	URL := viper.GetString("ZT.URL")
	NetworkID := viper.GetString("ZT.Network")
	suffix := viper.GetString("suffix")
	ztnetwork := ztapi.GetNetworkInfo(API, URL, NetworkID)
	log.Infof("Getting Members of Network: %s", ztnetwork.Config.Name)
	lst := ztapi.GetMemberList(API, URL, ztnetwork.ID)
	log.Infof("Got %d members", len(*lst))
	for _, n := range *lst {
		if n.Online {
			record := n.Name + "." + suffix + "."
			dnssrv.DNSDatabase[record] = dnssrv.Records{}
			ip6 := []net.IP{}
			ip4 := []net.IP{}
			if ztnetwork.Config.V6AssignMode.Sixplane {
				ip6 = append(ip6, n.Get6Plane())
			}
			if ztnetwork.Config.V6AssignMode.Rfc4193 {
				ip6 = append(ip6, n.GetRFC4193())
			}
			for _, a := range n.Config.IPAssignments {
				ip4 = append(ip4, net.ParseIP(a))
			}
			log.Infof("Updating %-15s IPv4: %-15s IPv6: %s", record, ip4, ip6)
			dnssrv.DNSDatabase[record] = dnssrv.Records{
				A:    ip4,
				AAAA: ip6,
			}
		}
	}
	return time.Now()
}
