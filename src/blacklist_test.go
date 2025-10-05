package main

import (
	"testing"

	"github.com/miekg/dns"
)

func assertWhitelisted(t *testing.T, domain string, expected bool) {
	result := isWhitelisted(domain)
	checkTestBool(t, expected, result)
}

func assertBlacklisted(t *testing.T, domain string, expected bool) {
	result := isBlacklisted(domain)
	checkTestBool(t, expected, result)
}

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

func TestWildcardWhitelisted(t *testing.T) {
	// We'll implement our test by temporarily changing the global state
	originalWhitelist := whitelistRecords
	whitelistRecords = []string{"*.google.com"}
	defer func() { whitelistRecords = originalWhitelist }()

	assertWhitelisted(t, "mail.google.com", true)
	assertWhitelisted(t, "google.com", true)
	assertWhitelisted(t, "deep.deep.google.com", true)
	assertWhitelisted(t, "google.org", false)
}

func TestWildcardBlacklisted(t *testing.T) {
	// Set configuration to blacklist everything
	originalConfig := ConfigInstance
	ConfigInstance = &Config{BlacklistEverything: false}
	defer func() { ConfigInstance = originalConfig }()

	// Test with specific blacklist entries containing wildcards
	originalBlacklist := blacklistRecords
	blacklistRecords = []string{"*.bad-domain.com"}
	defer func() { blacklistRecords = originalBlacklist }()
	assertBlacklisted(t, "bad-domain.com", true)
	assertBlacklisted(t, "under.bad-domain.com", true)
	assertBlacklisted(t, "under.deep.bad-domain.com", true)
	assertBlacklisted(t, "bing.com", false)
}
