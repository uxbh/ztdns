// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package cmd

import (
	"log"
	"net"
	"time"

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
		if viper.GetString("API") == "" {
			log.Fatal("No API key provided")
		}
		if viper.GetString("ZT.Network") == "" {
			log.Fatal("No Network ID Provided")
		}
		lastUpdate := updateDNS()
		req := make(chan bool)
		go dnssrv.Start(viper.GetInt("port"), viper.GetString("suffix"), req)
		for {
			<-req
			if time.Since(lastUpdate) > 30*time.Minute {
				log.Printf("DNSDatabase is stale. Refreshing.")
				lastUpdate = updateDNS()
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func updateDNS() time.Time {
	API := viper.GetString("API")
	URL := viper.GetString("ZT.URL")
	NetworkID := viper.GetString("ZT.Network")
	suffix := viper.GetString("suffix")
	ztnetwork := ztapi.GetNetworkInfo(API, URL, NetworkID)
	log.Printf("Got Network: %s\n", ztnetwork.Config.Name)
	lst := ztapi.GetMemberList(API, URL, ztnetwork.ID)
	log.Printf("Got %d members", len(*lst))
	for _, n := range *lst {
		if n.Online {
			var ip6 net.IP
			var ip4 net.IP
			switch {
			case ztnetwork.Config.V6AssignMode.Sixplane:
				ip6 = n.Get6Plane()
			case ztnetwork.Config.V6AssignMode.Rfc4193:
				ip6 = n.GetRFC4193()
			default:
				ip6 = nil
			}
			switch {
			case len(n.Config.IPAssignments) > 0:
				ip4 = net.ParseIP(n.Config.IPAssignments[0])
			default:
				ip4 = nil
			}
			//fmt.Println(n.Name, n.NetworkID, n.NodeID, ip4, ip6)
			record := n.Name + "." + suffix + "."
			log.Printf("Updating %-15s IPv4: %-15s IPv6: %s", record, ip4, ip6)
			dnssrv.DNSDatabase[record] = dnssrv.Records{
				A:    ip4,
				AAAA: ip6,
			}
		}
	}
	return time.Now()
}
