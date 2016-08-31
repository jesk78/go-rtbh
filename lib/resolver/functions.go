package resolver

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/orm"
	"math/rand"
	"net"
	"time"
)

const MAX_SAMPLES int = 100
const MAX_SLEEP_INTERVAL int64 = 5000

func (r *Resolver) RandomSelectSamples(num int) []string {
	var addr *orm.Address
	var all_nofqdn []string
	var sampling []string
	var sample string
	var sample_known bool
	var sampled string

	for _, addr = range orm.GetAddressesNoFqdn() {
		all_nofqdn = append(all_nofqdn, addr.Addr)
	}

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
		sample = all_nofqdn[rand.Intn(len(all_nofqdn))]
		sample_known = false

		for _, sampled = range sampling {
			if sampled == sample {
				sample_known = true
				break
			}
		}

		if !sample_known {
			sampling = append(sampling, sample)
		}
	}

	return sampling
}

func (r *Resolver) Lookup(addr string) string {
	var cached_addr string
	var cached_fqdn string
	var names []string
	var fqdn string
	var err error

	// First, check if we already have an entry in cache
	for cached_addr, cached_fqdn = range r.cache {
		if cached_addr == addr {
			return cached_fqdn
		}
	}

	// If we dont find anything, perform a DNS lookup
	if names, err = net.LookupAddr(addr); err != nil {
		Log.Warning("[DNSLookup]: Failed to lookup address for " + addr + ": " + err.Error())
		return ""
	}
	fqdn = names[0]

	r.mutex.Lock()
	r.cache[addr] = fqdn
	r.mutex.Unlock()

	return fqdn
}

func (r *Resolver) UnknownLookupRoutine() (err error) {
	var addr string
	var fqdn string
	var sampling []string
	var stop_loop bool
	var t_tick time.Time
	var t_now time.Time

	t_now = time.Now()
	t_tick = t_now.Add(time.Duration(rand.Int63n(MAX_SLEEP_INTERVAL)))

	Log.Debug(MYNAME + ": Starting UnknownLookupRoutine")
	stop_loop = false
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
						Log.Debug("Shutting down dnslookup")
						stop_loop = true
						continue
					}
				}
			}
		default:
			{
				if t_now = time.Now(); t_now.After(t_tick) {
					sampling = r.RandomSelectSamples(1)
					if len(sampling) == 1 {
						addr = sampling[0]
						if fqdn = r.Lookup(addr); fqdn != "" {
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

	return
}
