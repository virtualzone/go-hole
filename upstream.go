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

	c := new(dns.Client)
	in, _, err := c.Exchange(m1, GetConfig().UpstreamDNS[0])
	if err != nil {
		return nil, errors.New("failed to query upstream DNS server: " + err.Error())
	}
	return in.Answer, nil
}
