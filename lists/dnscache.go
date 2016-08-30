package lists

import (
	"net"
	"sync"
)

type dnsCache struct {
	Cache map[string]string
	Mutex sync.Mutex
}

var dns dnsCache = dnsCache{
	Cache: make(map[string]string),
}

func DNSLookup(addr string) string {
	var cached_addr string
	var cached_fqdn string
	var names []string
	var fqdn string
	var err error

	// First, check if we already have an entry in cache
	for cached_addr, cached_fqdn = range dns.Cache {
		if cached_addr == addr {
			Log.Info("Returning entry from DNS Cache")
			return cached_fqdn
		}
	}

	// If we dont find anything, perform a DNS lookup
	if names, err = net.LookupAddr(addr); err != nil {
		Log.Warning("[DNSLookup]: Failed to lookup address for " + addr + ": " + err.Error())
		return ""
	}
	fqdn = names[0]

	dns.Mutex.Lock()
	dns.Cache[addr] = fqdn
	dns.Mutex.Unlock()

	return fqdn
}
