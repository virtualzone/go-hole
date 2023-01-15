package main

import (
	"testing"

	"github.com/miekg/dns"
)

func TestLocalASuccess(t *testing.T) {
	res, err := queryLocal("service1.local.", dns.TypeA)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 2, len(res))
	aRecord1 := res[0].(*dns.A)
	aRecord2 := res[1].(*dns.A)
	checkTestString(t, "192.168.178.1", aRecord1.A.String())
	checkTestString(t, "192.168.179.1", aRecord2.A.String())
}

func TestLocalAAAASuccess(t *testing.T) {
	res, err := queryLocal("service1.local.", dns.TypeAAAA)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 1, len(res))
	aRecord1 := res[0].(*dns.AAAA)
	checkTestString(t, "fe80::9656:d028:8652:1111", aRecord1.AAAA.String())
}

func TestLocalNonExistent(t *testing.T) {
	res, err := queryLocal("nonexistentrecord.local.", dns.TypeA)
	checkTestBool(t, false, err == nil)
	checkTestBool(t, true, res == nil)
}
