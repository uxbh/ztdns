// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

// Package dnssrv implements a simple DNS server.
package dnssrv

import (
	"fmt"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

// Records contains the types of records the server will respond to.
type Records struct {
	A    []net.IP
	AAAA []net.IP
}

// DNSUpdate is the last time the DNSDatabase was updated.
var DNSUpdate = time.Time{}

// DNSDatabase is a map of hostnames to the records associated with it.
var DNSDatabase = map[string]Records{}

var queryChan chan string

// Start brings up a DNS server for the specified suffix on a given port.
func Start(iface string, port int, suffix string, req chan string) error {
	queryChan = req

	if port == 0 {
		port = 53
	}

	// attach request handler func
	dns.HandleFunc(suffix + ".", handleDNSRequest)

	for _, addr := range getIfaceAddrs(iface) {
		go func(suffix string, addr net.IP, port int) {
			var server *dns.Server
			if addr.To4().String() == addr.String() {
				log.Debugf("Creating IPv4 Server: %s:%d udp", addr, port)
				server = &dns.Server{
					Addr: fmt.Sprintf("%s:%d", addr, port),
					Net:  "udp",
				}
			} else {
				log.Debugf("Creating IPv6 Server: [%s]:%d udp6", addr, port)
				server = &dns.Server{
					Addr: fmt.Sprintf("[%s]:%d", addr, port),
					Net:  "udp6",
				}
			}
			log.Printf("Starting server for %s on %s", suffix, server.Addr)
			err := server.ListenAndServe()
			if err != nil {
				log.Fatalf("failed to start DNS server: %s", err.Error())
			}
			defer server.Shutdown()
		}(suffix, addr, port)
	}
	return nil
}

func getIfaceAddrs(iface string) []net.IP {
	if iface != "" {
		retaddrs := []net.IP{}
		netint, err := net.InterfaceByName(iface)
		if err != nil {
			log.Fatalf("Could not get interface: %s\n", err.Error())
		}
		addrs, err := netint.Addrs()
		if err != nil {
			log.Fatalf("Could not get addresses: %s\n", err.Error())
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if !ip.IsLinkLocalUnicast() {
				log.Debugf("Found address: %s", ip.String())
				retaddrs = append(retaddrs, ip)
			}
		}
		return retaddrs
	}
	return []net.IP{net.IPv4zero}
}

// handleDNSRequest routes an incoming DNS request to a parser.
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

// parseQuery reads and creates an answer to a DNS query.
func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		queryChan <- q.Name
		if rec, ok := DNSDatabase[q.Name]; ok {
			switch q.Qtype {
			case dns.TypeA:
				for _, ip := range rec.A {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip.String()))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			case dns.TypeAAAA:
				for _, ip := range rec.AAAA {
					rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip.String()))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
		}
	}
}
