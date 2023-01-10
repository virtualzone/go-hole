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
		found := false
		if q.Qtype == dns.TypeA {
			arr, err := queryLocal(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				logQueryResult(source, name, q.Qtype, "resolved as local address")
			}
		}
		if !found {
			arr, err := queryBlacklist(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				logQueryResult(source, name, q.Qtype, "resolved as blacklisted name")
			}
		}
		if !found {
			arr, err := queryUpstream(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				logQueryResult(source, name, q.Qtype, "resolved via upstream")
			}
		}
		if !found {
			logQueryResult(source, name, q.Qtype, "did not resolve")
		}
	}
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
