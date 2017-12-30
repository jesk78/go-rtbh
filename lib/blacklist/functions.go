package blacklist

import (
	"time"

	"fmt"

	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/resolver"
)

func (bl *Blacklist) Add(event events.RTBHEvent) error {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	// Set fqdn to not-yet-lookedup so it will be picked up by the Resolver
	fqdn := resolver.FQDN_TO_LOOKUP

	addr, err := orm.UpdateAddress(event.Address, fqdn)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}
	if addr.Addr == "" {
		return fmt.Errorf("Blacklist.Add: address is empty")
	}

	reason, err := orm.UpdateReason(event.Reason)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}
	if reason.Reason == "" {
		return fmt.Errorf("Blacklist.Add: reason is empty")
	}

	entry, err := orm.GetBlacklistEntry(event.Address)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}

	expireOn, err := time.ParseDuration(event.ExpireIn)
	if err != nil {
		return fmt.Errorf("Blacklist.Add time.ParseDuration: %v", err)
	}

	entry = &orm.Blacklist{
		AddrId:   addr.Id,
		ReasonId: reason.Id,
		AddedAt:  event.AddedAt,
		ExpireOn: time.Now().Add(expireOn),
	}

	err = entry.Save(addr.Addr)
	if err != nil {
		return fmt.Errorf("Blacklist.Add: %v", err)
	}

	log.Debugf("Blacklist.Add: Adding BGP route")
	bl.bgp.AddRoute(addr.Addr)

	return nil
}

func (bl *Blacklist) Remove(addr string) error {
	entry, err := orm.GetBlacklistEntry(addr)
	if err != nil {
		return fmt.Errorf("Blacklist.Remove: %v", err)
	}

	err = entry.Remove()
	if err != nil {
		return fmt.Errorf("Blacklist.Remove: %v", err)
	}

	bl.bgp.RemoveRoute(addr)

	return nil
}

func (bl *Blacklist) Listed(addr string) bool {
	entry, err := orm.GetBlacklistEntry(addr)
	return entry != nil && err == nil
}

func (bl *Blacklist) ReapExpiredEntries() error {
	t_now := time.Now()

	for _, entry := range orm.GetBlacklistEntries() {
		addr, err := orm.GetAddressById(entry.AddrId)
		if err != nil {
			return fmt.Errorf("Blacklist.ReapExpiredEntries: %v", err)
		}
		if addr.Addr == "" {
			return fmt.Errorf("Blacklist.ReapExpiredEntries: Failed to locate address record for object id %d", entry.Id)
		}

		if t_now.After(entry.ExpireOn) {
			err = entry.Remove()
			if err != nil {
				log.Warningf("Blacklist.ReapExpiredEntries: %v", err)
				continue
			}
			log.Debugf("Blacklist.ReapExpiredEntries: %s expired from blacklist", addr.Addr)
		}
	}

	return nil
}

func (bl *Blacklist) GetById(id int64) (*events.APIEvent, error) {
	for _, entry := range orm.GetBlacklistEntries() {
		if entry.Id == id {
			addr, err := orm.GetAddressById(entry.AddrId)
			if err != nil {
				return nil, fmt.Errorf("Blacklist.GetById: %v", err)
			}
			if addr.Addr == "" {
				return nil, fmt.Errorf("Blacklist.GetById: Did not find ip address for blacklist entry for id %d", id)
			}

			reason, err := orm.GetReasonById(entry.ReasonId)
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
	for _, entry := range orm.GetBlacklistEntries() {
		addr, err := orm.GetAddressById(entry.AddrId)
		if err != nil {
			return nil, fmt.Errorf("Blacklist.GetAll: %v", err)
		}
		if addr.Addr == "" {
			return nil, fmt.Errorf("Blacklist.GetAll: Did not find ip address for blacklist entry for object id %d", entry.Id)
		}

		reason, err := orm.GetReasonById(entry.ReasonId)
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
