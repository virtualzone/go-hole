package main

import (
	"errors"

	"github.com/miekg/dns"
)

func queryUpstream(name string, qtype uint16) ([]dns.RR, error) {
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
