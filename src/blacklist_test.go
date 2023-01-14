package main

import (
	"testing"

	"github.com/miekg/dns"
)

func TestBlacklistSuccess(t *testing.T) {
	res, err := queryBlacklist("googleads.g.doubleclick.net.", dns.TypeA)
	checkTestBool(t, true, err == nil)
	checkTestBool(t, false, res == nil)
	checkTestInt(t, 0, len(res))
}

func TestBlacklistNonExistent(t *testing.T) {
	res, err := queryBlacklist("www.apple.com.", dns.TypeA)
	checkTestBool(t, false, err == nil)
	checkTestBool(t, true, res == nil)
}
