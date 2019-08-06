// Copyright © 2017 uxbh
// This file is part of github.com/hatemosphere/ztdns.

package cmd

import (
	"fmt"
	"net"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/hatemosphere/ztdns/dnssrv"
	"github.com/hatemosphere/ztdns/ztapi"
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

	rrDNSPatterns := make(map[string]*regexp.Regexp)
	rrDNSRecords := make(map[string][]dnssrv.Records)

	for re, host := range viper.GetStringMapString("RoundRobin") {
		rrDNSPatterns[host] = regexp.MustCompile(re)
		log.Debugf("Creating match '%s' for %s host", re, host)
	}

	// Get all configured networks:
	for domain, id := range viper.GetStringMapString("Networks") {
		// Get ZeroTier Network info
		ztnetwork, err := ztapi.GetNetworkInfo(API, URL, id)
		if err != nil {
			log.Fatalf("Unable to update DNS entries: %s", err.Error())
		}

		// Get list of members in network
		log.Infof("Getting Members of Network: %s (%s)", ztnetwork.Config.Name, domain)
		lst, err := ztapi.GetMemberList(API, URL, ztnetwork.ID)
		if err != nil {
			log.Fatalf("Unable to update DNS entries: %s", err.Error())
		}
		log.Infof("Got %d members", len(*lst))

		for _, n := range *lst {
			// For all online members
			if n.Online {
				// Clear current DNS records
				record := n.Name + "." + domain + "." + suffix + "."
				dnssrv.DNSDatabase[record] = dnssrv.Records{}
				ip6 := []net.IP{}
				ip4 := []net.IP{}
				// Get 6Plane address if network has it enabled
				if ztnetwork.Config.V6AssignMode.Sixplane {
					ip6 = append(ip6, n.Get6Plane())
				}
				// Get RFC4193 address if network has it enabled
				if ztnetwork.Config.V6AssignMode.Rfc4193 {
					ip6 = append(ip6, n.GetRFC4193())
				}

				// Get the rest of the address assigned to the member
				for _, a := range n.Config.IPAssignments {
					ip4 = append(ip4, net.ParseIP(a))
				}

				dnsRecord := dnssrv.Records{
					A:    ip4,
					AAAA: ip6,
				}

				// Add the record to the database
				log.Infof("Updating %-15s IPv4: %-15s IPv6: %s", record, ip4, ip6)
				dnssrv.DNSDatabase[record] = dnsRecord

				// Finding matches for RoundRobin dns
				for host, re := range rrDNSPatterns {
					log.Debugf("Checking matches for %s host", host)
					if match := re.FindStringSubmatch(n.Name); match != nil {
						// prefix := fmt.Sprintf(host, iface(match[1:]))
						rrRecord := host + "." + domain + "." + suffix + "."

						log.Infof("Adding ips to RR record %-15s IPv4: %-15s IPv6: %s, from host %s", rrRecord, ip4, ip6, n.Name)
						rrDNSRecords[rrRecord] = append(rrDNSRecords[rrRecord], dnsRecord)
					}
				}
			}
		}

		for rrRecord, dnsRecords := range rrDNSRecords {
			rrRecordIps := dnssrv.Records{}
			for _, ips := range dnsRecords {
				rrRecordIps.A = append(rrRecordIps.A, ips.A...)
				rrRecordIps.AAAA = append(rrRecordIps.AAAA, ips.AAAA...)
			}

			log.Infof("Updating %-15s IPv4: %-15s IPv6: %s", rrRecord, rrRecordIps.A, rrRecordIps.AAAA)
			dnssrv.DNSDatabase[rrRecord] = rrRecordIps
		}
	}

	// Return the current update time
	return time.Now()
}

// Convert slice of string to interface for fmt
func iface(list []string) []interface{} {
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}
	return vals
}
