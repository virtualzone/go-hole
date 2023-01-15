package main

import (
	"errors"
	"log"

	"github.com/miekg/dns"
)

func queryUpstream(name string, qtype uint16) ([]dns.RR, error) {
	// Check cache first
	res, err := GetUpstreamCache().Get(name, qtype)
	if err == nil {
		// Record found in cache
		log.Printf("query for %s %s resolved via cache\n", getQueryTypeText(qtype), name)
		if len(res) == 0 {
			return nil, errors.New("record not found via upstream DNS server")
		}
		return res, nil
	}

	// If not cached, perform actual upstream query
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{
		Name:   name,
		Qtype:  qtype,
		Qclass: dns.ClassINET,
	}

	for _, server := range GetConfig().UpstreamDNS {
		in, err := doUpstreamQuery(m1, server)
		if err == nil {
			GetUpstreamCache().Set(name, qtype, in.Answer)
			if len(in.Answer) == 0 {
				return nil, errors.New("record not found via upstream DNS server")
			}
			return in.Answer, nil
		}
	}

	return nil, errors.New("could not resolve query via any upstream DNS server")
}

func doUpstreamQuery(m *dns.Msg, address string) (*dns.Msg, error) {
	c := new(dns.Client)
	in, _, err := c.Exchange(m, address)
	return in, err
}
