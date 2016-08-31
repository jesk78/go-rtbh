package whitelist

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/orm"
	"net"
)

func (wl *Whitelist) Add(entry events.RTBHWhiteEntry) (err error) {
	var (
		addr   *orm.Address
		wentry orm.Whitelist
		names  []string
		fqdn   string
	)

	if names, err = net.LookupAddr(entry.Address); err != nil {
		Log.Warning(MYNAME + ": Failed to lookup fqdn for " + entry.Address)
		fqdn = "unknown"
	} else {
		fqdn = names[0]
	}

	if len(names) > 1 {
		Log.Warning(MYNAME + ": Multiple hosts found for " + entry.Address + " using " + fqdn)
	}

	if addr = orm.UpdateAddress(entry.Address, fqdn); addr.Addr == "" {
		return
	}

	wentry = orm.Whitelist{
		AddrId:      addr.Id,
		Description: entry.Description,
	}
	if ok := wentry.Save(); !ok {
		return
	}

	wl.bgp.RemoveRoute(entry.Address)

	wl.cache.Add(entry.Address, entry)

	return
}

func (wl *Whitelist) Remove(addr string) (err error) {
	var entry *orm.Whitelist

	if entry = orm.GetWhitelistEntry(addr); entry == nil {
		err = errors.New(MYNAME + ": Failed to retrieve address")
		return
	}

	wl.cache.Remove(addr)

	if ok := entry.Remove(); !ok {
		err = errors.New(MYNAME + ": Failed to remove entry")
	}

	return
}

func (wl *Whitelist) Listed(addr string) bool {
	return wl.cache.Has(addr)
}
