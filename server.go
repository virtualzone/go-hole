package main

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		name := strings.ToLower(q.Name)
		log.Printf("Query for %s\n", name)
		found := false
		if q.Qtype == dns.TypeA {
			log.Println("Checking local...")
			arr, err := queryLocal(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
			}
		}
		if !found {
			log.Println("Checking blacklist...")
			arr, err := queryBlacklist(name, q.Qtype)
			if err == nil {
				found = true
				m.Answer = append(m.Answer, arr...)
			}
		}
		if !found {
			log.Println("Checking upstream...")
			arr, err := queryUpstream(name, q.Qtype)
			if err == nil {
				m.Answer = append(m.Answer, arr...)
			}
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
