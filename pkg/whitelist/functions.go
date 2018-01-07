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
		wl.log.Warningf("Whitelist.Add: Failed to lookup fqdn for " + entry.Address)
	} else {
		fqdn = names[0]
	}

	if len(names) > 1 {
		wl.log.Warningf("Whitelist.Add: Multiple hosts found for " + entry.Address + " using " + fqdn)
	}

	addr, err := wl.orm.UpdateAddress(entry.Address, fqdn)
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
	entry, err := wl.orm.GetWhitelistEntry(addr)
	if err != nil {
		return fmt.Errorf("Whitelist.Remove: %v", err)
	}

	err = entry.Remove()
	if err != nil {
		return fmt.Errorf("Whitelist.Remove: %v", err)
	}

	return nil
}

func (wl *Whitelist) Update(data events.RTBHWhiteEntry) error {
	entry, err := wl.orm.GetWhitelistEntry(data.Address)
	if err != nil {
		return fmt.Errorf("Whitelist.Update: %v", err)
	}

	entry.Description = data.Description

	err = entry.Update()
	if err != nil {
		return fmt.Errorf("Whitelist.Update entry.Update: %v", err)
	}

	return nil
}

func (wl *Whitelist) Listed(addr string) bool {
	entry, err := wl.orm.GetWhitelistEntry(addr)
	return err == nil && entry.AddrId != 0
}

func (wl *Whitelist) GetAll() ([]*events.RTBHWhiteEntry, error) {
	entries := []*events.RTBHWhiteEntry{}

	wlEntries, err := wl.orm.GetWhitelistEntries()
	if err != nil {
		return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
	}

	for _, entry := range wlEntries {
		addr, err := wl.orm.GetAddressById(entry.AddrId)
		if err != nil {
			return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
		}
		if addr.Addr == "" {
			return nil, fmt.Errorf("Blacklist.GetAll: Did not find ip address for blacklist entry for object id %d", entry.Id)
		}

		event := &events.RTBHWhiteEntry{
			Id:          entry.Id,
			Address:     addr.Addr,
			Description: entry.Description,
		}

		entries = append(entries, event)
	}

	return entries, nil
}
