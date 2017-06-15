// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

package dnssrv

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

// Records contains the types of records the server will respond to.
type Records struct {
	A    net.IP
	AAAA net.IP
}

// DNSDatabase is a map of hostnames to the records associated with it.
var DNSDatabase = map[string]Records{}

// Start brings up a DNS server for the specified suffix on a given port.
func Start(port int, suffix string) error {
	dns.HandleFunc(suffix, handleDNSRequest)

	server := &dns.Server{
		Addr: fmt.Sprintf(":%d", port),
		Net:  "udp",
	}
	fmt.Printf("Starting server for %s on %d\n", suffix, port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start DNS server: %s", err.Error())
	}
	defer server.Shutdown()
	return nil
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
		if rec, ok := DNSDatabase[q.Name]; ok {
			switch q.Qtype {
			case dns.TypeA:
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, rec.A.String()))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			case dns.TypeAAAA:
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, rec.AAAA.String()))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}
