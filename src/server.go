package main

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func parseQuery(source net.Addr, m *dns.Msg) {
	for _, q := range m.Question {
		name := strings.ToLower(q.Name)
		processDnsQuery(name, q.Qtype, source, m)
	}
}

func processDnsQuery(name string, qtype uint16, source net.Addr, m *dns.Msg) {
	if qtype == dns.TypeA {
		arr, err := queryLocal(name, qtype)
		if err == nil {
			m.Answer = append(m.Answer, arr...)
			m.Rcode = dns.RcodeSuccess
			logQueryResult(source, name, qtype, "resolved as local address")
			return
		}
	}
	arr, err := queryBlacklist(name, qtype)
	if err == nil {
		m.Answer = append(m.Answer, arr...)
		m.Rcode = dns.RcodeNameError
		logQueryResult(source, name, qtype, "resolved as blacklisted name")
		return
	}
	arr, err = queryUpstream(name, qtype)
	if err == nil {
		m.Answer = append(m.Answer, arr...)
		m.Rcode = dns.RcodeSuccess
		logQueryResult(source, name, qtype, "resolved via upstream")
		return
	}
	m.Rcode = dns.RcodeNameError
	logQueryResult(source, name, qtype, "did not resolve")
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(w.RemoteAddr(), m)
	}

	w.WriteMsg(m)
	w.Close()
}

func listenAndServe() {
	dns.HandleFunc(".", handleDnsRequest)

	server := &dns.Server{
		Addr: GetConfig().ListenAddr,
		Net:  "udp",
	}
	log.Printf("Starting at %s\n", GetConfig().ListenAddr)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
