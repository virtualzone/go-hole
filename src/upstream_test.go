package main

import (
	"testing"

	"github.com/miekg/dns"
)

func TestUpstreamSuccess(t *testing.T) {
	res, err := queryUpstream("dns.google.", dns.TypeA)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 2, len(res))
	aRecord1 := res[0].(*dns.A)
	aRecord2 := res[1].(*dns.A)
	checkTestBool(t, true, aRecord1.A.String() == "8.8.8.8" || aRecord1.A.String() == "8.8.4.4")
	checkTestBool(t, true, aRecord2.A.String() == "8.8.8.8" || aRecord2.A.String() == "8.8.4.4")
	checkTestBool(t, true, aRecord1.A.String() != aRecord2.A.String())
}

func TestUpstreamNonExistent(t *testing.T) {
	res, err := queryUpstream("nonexistentrecord.virtualzone.de.", dns.TypeA)
	checkTestBool(t, false, err == nil)
	checkTestBool(t, true, res == nil)
}

func TestUpstreamCname(t *testing.T) {
	res, err := queryUpstream("iadsdk.apple.com.", dns.TypeCNAME)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 1, len(res))
	cnameRecord1 := res[0].(*dns.CNAME)
	checkTestString(t, "iadsdk.apple.com.akadns.net.", cnameRecord1.Target)
}
