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
		res, errCode := processDnsQuery(name, q.Qtype, source)
		m.Answer = append(m.Answer, res...)
		m.Rcode = errCode
	}
}

func processDnsQuery(name string, qtype uint16, source net.Addr) ([]dns.RR, int) {
	arr, err := queryLocal(name, qtype)
	if err == nil {
		logQueryResult(source, name, qtype, "resolved as local address")
		return arr, dns.RcodeSuccess
	}
	arr, err = queryBlacklist(name, qtype)
	if err == nil {
		logQueryResult(source, name, qtype, "resolved as blacklisted name")
		return arr, dns.RcodeNameError
	}
	arr, err = queryUpstream(name, qtype)
	if err == nil {
		logQueryResult(source, name, qtype, "resolved via upstream")
		return arr, dns.RcodeSuccess
	}
	logQueryResult(source, name, qtype, "did not resolve")
	return []dns.RR{}, dns.RcodeNameError
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
