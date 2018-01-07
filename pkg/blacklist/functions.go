package blacklist

import (
	"time"

	"fmt"

	"strings"

	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/orm"
	"github.com/r3boot/go-rtbh/pkg/resolver"
)

func (bl *Blacklist) Add(event events.RTBHEvent) error {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.log.Debugf("event: %v", event.Reason)

	// Set fqdn to not-yet-lookedup so it will be picked up by the Resolver
	fqdn := resolver.FQDN_TO_LOOKUP

	addr, err := bl.orm.UpdateAddress(event.Address, fqdn)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}
	if addr == nil {
		return fmt.Errorf("Blacklist.Add: address is empty")
	}

	reason, err := bl.orm.UpdateReason(event.Reason)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}

	if reason.Reason == "" {
		return fmt.Errorf("Blacklist.Add: reason is empty")
	}

	expireOn, err := time.ParseDuration(event.ExpireIn)
	if err != nil {
		return fmt.Errorf("Blacklist.Add time.ParseDuration: %v", err)
	}

	entry := &orm.Blacklist{
		AddrId:   addr.Id,
		ReasonId: reason.Id,
		AddedAt:  event.AddedAt,
		ExpireOn: time.Now().Add(expireOn),
	}

	err = entry.Save(addr.Addr)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}

	bl.log.Debugf("Blacklist.Add: Adding BGP route")
	bl.bgp.AddRoute(addr.Addr)

	return nil
}

func (bl *Blacklist) Remove(addr string) error {
	entry, err := bl.orm.GetBlacklistEntry(addr)
	if err != nil {
		return fmt.Errorf("Blacklist.Remove: %v", err)
	}

	err = entry.Remove()
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("Blacklist.Remove: %v", err)
		}
	}

	bl.bgp.RemoveRoute(addr)

	return nil
}

func (bl *Blacklist) Listed(addr string) bool {
	entry, err := bl.orm.GetBlacklistEntry(addr)
	return entry != nil && err == nil
}

func (bl *Blacklist) ReapExpiredEntries() error {
	t_now := time.Now()

	entries, err := bl.orm.GetBlacklistEntries()
	if err != nil {
		return fmt.Errorf("Blacklist.ReapExpiredEntries: %v", err)
	}

	for _, entry := range entries {
		addr, err := bl.orm.GetAddressById(entry.AddrId)
		if err != nil {
			return fmt.Errorf("Blacklist.ReapExpiredEntries: %v", err)
		}
		if addr.Addr == "" {
			return fmt.Errorf("Blacklist.ReapExpiredEntries: Failed to locate address record for object id %d", entry.Id)
		}

		if t_now.After(entry.ExpireOn) {
			err = entry.Remove()
			if err != nil {
				bl.log.Warningf("Blacklist.ReapExpiredEntries: %v", err)
				continue
			}
			bl.log.Debugf("Blacklist.ReapExpiredEntries: %s expired from blacklist", addr.Addr)
		}
	}

	return nil
}

func (bl *Blacklist) GetById(id int64) (*events.APIEvent, error) {
	entries, err := bl.orm.GetBlacklistEntries()
	if err != nil {
		return nil, fmt.Errorf("Blacklist.GetById: %v", err)
	}

	for _, entry := range entries {
		if entry.Id == id {
			addr, err := bl.orm.GetAddressById(entry.AddrId)
			if err != nil {
				return nil, fmt.Errorf("Blacklist.GetById: %v", err)
			}
			if addr.Addr == "" {
				return nil, fmt.Errorf("Blacklist.GetById: Did not find ip address for blacklist entry for id %d", id)
			}

			reason, err := bl.orm.GetReasonById(entry.ReasonId)
			if err != nil {
				return nil, fmt.Errorf("Blacklist.GetById: %v", err)
			}
			if reason.Reason == "" {
				return nil, fmt.Errorf("Blacklist.GetById: Did not find reason for blacklist entry for id %d", id)
			}

			return &events.APIEvent{
				Id:       entry.Id,
				Address:  addr.Addr,
				Reason:   reason.Reason,
				AddedAt:  entry.AddedAt,
				ExpireOn: entry.ExpireOn,
			}, nil
		}
	}

	return nil, fmt.Errorf("Blacklist.GetById: No such id")
}

func (bl *Blacklist) GetAll() ([]*events.APIEvent, error) {
	entries := []*events.APIEvent{}

	blEntries, err := bl.orm.GetBlacklistEntries()
	if err != nil {
		return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
	}

	for _, entry := range blEntries {
		addr, err := bl.orm.GetAddressById(entry.AddrId)
		if err != nil {
			return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
		}
		if addr.Addr == "" {
			return nil, fmt.Errorf("Blacklist.GetAll: Did not find ip address for blacklist entry for object id %d", entry.Id)
		}

		reason, err := bl.orm.GetReasonById(entry.ReasonId)
		if err != nil {
			return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
		}
		if reason.Reason == "" {
			return nil, fmt.Errorf("Blacklist.GetAll: Did not find reason for blacklist entry for %s", addr.Addr)
		}

		event := &events.APIEvent{
			Id:       entry.Id,
			Address:  addr.Addr,
			Reason:   reason.Reason,
			AddedAt:  entry.AddedAt,
			ExpireOn: entry.ExpireOn,
		}

		entries = append(entries, event)
	}

	return entries, nil
}
