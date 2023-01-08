package main

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		name := strings.ToLower(q.Name)
		found := false
		if q.Qtype == dns.TypeA {
			arr, err := queryLocal(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				log.Printf("Query for %s resolved as local address\n", name)
			}
		}
		if !found {
			arr, err := queryBlacklist(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				log.Printf("Query for %s resolved as blacklisted name\n", name)
			}
		}
		if !found {
			arr, err := queryUpstream(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
				log.Printf("Query for %s resolved via upstream\n", name)
			}
		}
		if !found {
			log.Printf("Query for %s did not resolve\n", name)
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
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
