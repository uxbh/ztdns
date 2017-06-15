// Copyright © 2017 uxbh
// This file is part of gitlab.com/uxbh/ztdns.

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
	A    net.IP
	AAAA net.IP
}

// DNSUpdate is the last time the DNSDatabase was updated.
var DNSUpdate = time.Time{}

// DNSDatabase is a map of hostnames to the records associated with it.
var DNSDatabase = map[string]Records{}

var queryChan chan bool

// Start brings up a DNS server for the specified suffix on a given port.
func Start(port int, suffix string, req chan bool) error {
	queryChan = req

	if port == 0 {
		port = 53
	}

	if suffix == "" {
		log.Fatal("No DNS Suffix provided.")
	}

	dns.HandleFunc(suffix, handleDNSRequest)

	server := &dns.Server{
		Addr: fmt.Sprintf(":%d", port),
		Net:  "udp",
	}
	log.Printf("Starting server for %s on %d", suffix, port)
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
	queryChan <- true
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
