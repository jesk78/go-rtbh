package whitelist

import (
	"net"

	"fmt"

	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func (wl *Whitelist) Add(entry events.RTBHWhiteEntry) error {
	var (
		wentry orm.Whitelist
		names  []string
	)

	fqdn := "UNRESOLVED"
	names, err := net.LookupAddr(entry.Address)
	if err != nil {
		log.Warningf("Whitelist.Add: Failed to lookup fqdn for " + entry.Address)
	} else {
		fqdn = names[0]
	}

	if len(names) > 1 {
		log.Warningf("Whitelist.Add: Multiple hosts found for " + entry.Address + " using " + fqdn)
	}

	addr, err := orm.UpdateAddress(entry.Address, fqdn)
	if err != nil {
		return fmt.Errorf("Whitelist.Add: %v", err)
	}
	if addr.Addr == "" {
		return fmt.Errorf("Whitelist.Add: Address is empty")
	}

	wentry = orm.Whitelist{
		AddrId:      addr.Id,
		Description: entry.Description,
	}
	err = wentry.Save()
	if err != nil {
		return fmt.Errorf("Whitelist.Add: %v", err)
	}

	wl.bgp.RemoveRoute(entry.Address)

	return nil
}

func (wl *Whitelist) Remove(addr string) error {
	entry, err := orm.GetWhitelistEntry(addr)
	if err != nil {
		return fmt.Errorf("Whitelist.Remove: %v", err)
	}

	err = entry.Remove()
	if err != nil {
		return fmt.Errorf("Whitelist.Remove: %v", err)
	}

	return nil
}

func (wl *Whitelist) Listed(addr string) bool {
	entry, err := orm.GetWhitelistEntry(addr)
	return entry.AddrId != 0 && err != nil
}
