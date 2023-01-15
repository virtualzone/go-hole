package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/miekg/dns"
)

type UpstreamCache struct {
	Cache *ttlcache.Cache[string, []dns.RR]
}

var UpstreamCacheInstance *UpstreamCache = &UpstreamCache{}

func GetUpstreamCache() *UpstreamCache {
	return UpstreamCacheInstance
}

func (c *UpstreamCache) Init() {
	c.Cache = ttlcache.New(
		ttlcache.WithDisableTouchOnHit[string, []dns.RR](),
	)
	go c.Cache.Start()
}

func (c *UpstreamCache) Set(name string, qtype uint16, rr []dns.RR) {
	ttl := time.Duration(c.getMinTtl(rr)) * time.Second
	c.Cache.Set(c.getKey(name, qtype), rr, ttl)
}

func (c *UpstreamCache) Get(name string, qtype uint16) ([]dns.RR, error) {
	res := c.Cache.Get(c.getKey(name, qtype))
	if res == nil || res.IsExpired() {
		return nil, errors.New("record not found in cache")
	}
	return res.Value(), nil
}

func (c *UpstreamCache) Clear() {
	c.Cache.DeleteAll()
}

func (c *UpstreamCache) getKey(name string, qtype uint16) string {
	return name + "_" + strconv.Itoa(int(qtype))
}

func (c *UpstreamCache) getMinTtl(rr []dns.RR) uint32 {
	var res uint32 = 1800 // Default: 30 minutes
	for _, record := range rr {
		if record.Header().Ttl < res {
			res = record.Header().Ttl
		}
	}
	return res
}
