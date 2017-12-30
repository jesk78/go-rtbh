package resolver

import (
	"math/rand"
	"net"
	"time"

	"fmt"

	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/orm"
)

func (r *Resolver) RandomSelectSamples(num int) ([]string, error) {
	entries, err := orm.GetAddressesNoFqdn()
	if err != nil {
		return nil, fmt.Errorf("Resolver.RandomSelectSamples: %v", err)
	}

	all_nofqdn := []string{}
	for _, addr := range entries {
		all_nofqdn = append(all_nofqdn, addr.Addr)
	}

	sampling := []string{}
	for {
		// Break if we have enough entries
		if len(sampling) >= MAX_SAMPLES {
			break
		}

		// Break if the amount of samples equals the input queue
		if len(all_nofqdn) == len(sampling) {
			break
		}

		// Pick a sample which is not yet sampled
		sample := all_nofqdn[rand.Intn(len(all_nofqdn))]
		sample_known := false

		for _, sampled := range sampling {
			if sampled == sample {
				sample_known = true
				break
			}
		}

		if !sample_known {
			sampling = append(sampling, sample)
		}
	}

	return sampling, nil
}

func (r *Resolver) Lookup(addr string) (string, error) {
	var cached_addr string
	var cached_fqdn string
	var names []string
	var fqdn string
	var err error

	// First, check if we already have an entry in cache
	for cached_addr, cached_fqdn = range r.cache {
		if cached_addr == addr {
			return cached_fqdn, nil
		}
	}

	// If we dont find anything, perform a DNS lookup
	if names, err = net.LookupAddr(addr); err != nil {
		return "", fmt.Errorf("Resolver.Lookup net.LookupAddr: %v", err)
	}
	fqdn = names[0]

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cache[addr] = fqdn

	return fqdn, nil
}

func (r *Resolver) UnknownLookupRoutine() error {
	t_now := time.Now()
	t_tick := t_now.Add(time.Duration(rand.Int63n(MAX_SLEEP_INTERVAL)))

	log.Debugf("Resolver: Starting UnknownLookupRoutine")
	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case cmd := <-r.Control:
			{
				switch cmd {
				case config.CTL_SHUTDOWN:
					{
						log.Debugf("Resolver: Shutting down UnknownLookupRoutine")
						stop_loop = true
						continue
					}
				}
			}
		default:
			{
				if t_now = time.Now(); t_now.After(t_tick) {
					sampling, err := r.RandomSelectSamples(1)
					if err != nil {
						log.Warningf("Resolver.UnknownLookupRoutine: %v", err)
						continue
					}

					if len(sampling) == 1 {
						addr := sampling[0]
						fqdn, err := r.Lookup(addr)
						if err != nil {
							log.Warningf("Resolver.UnknownLookupRoutine: %v", err)
							continue
						}
						if fqdn != "" {
							orm.UpdateAddress(addr, fqdn)
						}
					}
					t_tick = t_now.Add(time.Duration(rand.Int63n(MAX_SLEEP_INTERVAL)))
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}

	r.Done <- true

	return nil
}
