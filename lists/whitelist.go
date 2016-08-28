package lists

import (
	"errors"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/orm"
	"github.com/r3boot/go-rtbh/proto"
	"net"
)

type Whitelist struct {
}

func (obj Whitelist) Add(entry events.RTBHWhiteEntry) (err error) {
	var (
		addr   orm.Address
		wentry orm.Whitelist
		names  []string
		fqdn   string
	)

	if names, err = net.LookupAddr(entry.Address); err != nil {
		Log.Warning("[Whitelist]: Failed to lookup fqdn for " + entry.Address)
	}

	fqdn = names[0]
	if len(names) > 1 {
		Log.Warning("[Whitelist.Add]: Multiple hosts found for " + entry.Address + " using " + fqdn)
	}

	if addr = orm.UpdateAddress(entry.Address, fqdn); addr.Addr == "" {
		return
	}

	wentry = orm.Whitelist{
		Address:     &addr,
		Description: entry.Description,
	}
	if ok := wentry.Save(); !ok {
		return
	}

	proto.RemoveBGPRoute(entry.Address)

	return
}

func (obj Whitelist) Remove(addr string) (err error) {
	var entry orm.Whitelist

	if entry = orm.GetWhitelistEntry(addr); entry.Address.Addr == "" {
		err = errors.New("[Whitelist.Remove]: Failed to retrieve address")
		return
	}

	if ok := entry.Remove(); !ok {
		err = errors.New("[Whitelist.Remove]: Failed to remove entry")
	}

	return
}

func (obj Whitelist) Listed(addr string) bool {
	var entry orm.Whitelist

	entry = orm.GetWhitelistEntry(addr)
	return entry.Address != nil && entry.Address.Addr == addr
}
