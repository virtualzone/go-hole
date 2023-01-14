package main

import (
	"testing"

	"github.com/miekg/dns"
)

func TestLocalSuccess(t *testing.T) {
	res, err := queryLocal("service1.local.", dns.TypeA)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 1, len(res))
	aRecord1 := res[0].(*dns.A)
	checkTestBool(t, true, aRecord1.A.String() == "192.168.178.1")
}

func TestLocalNonExistent(t *testing.T) {
	res, err := queryLocal("nonexistentrecord.local.", dns.TypeA)
	checkTestBool(t, false, err == nil)
	checkTestBool(t, true, res == nil)
}
