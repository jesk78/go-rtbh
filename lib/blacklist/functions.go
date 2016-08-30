package blacklist

import (
	"errors"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/proto"
)

func (bl *Blacklist) Add(event events.RTBHEvent) (err error) {
	var (
		addr     orm.Address
		duration orm.Duration
		entry    orm.Blacklist
		fqdn     string
		reason   orm.Reason
	)

	// Set fqdn to not-yet-lookedup so it will be picked up by the Resolver
	fqdn = resolver.FQDN_TO_LOOKUP

	if addr = orm.UpdateAddress(event.Address, fqdn); addr.Addr == "" {
		err = errors.New(MYNAME + ": Failed to update Address record")
		return
	}

	if duration = orm.UpdateDuration(event.ExpireIn); duration.Duration == "" {
		err = errors.New(MYNAME + ": Failed to update Duration record")
		return
	}

	if reason = orm.UpdateReason(event.Reason); reason.Reason == "" {
		err = errors.New(MYNAME + ": Failed to update Reason record")
		return
	}

	entry = orm.Blacklist{
		AddressId:  addr.Id,
		ReasonId:   reason.Id,
		AddedAt:    event.AddedAt,
		DurationId: duration.Id,
	}
	if ok := entry.Save(); !ok {
		return
	}

	Log.Debug("Adding BGP route")
	proto.AddBGPRoute(addr.Addr)

	bl.cache.Add(addr.Addr, entry)

	return
}

func (bl *Blacklist) Remove(addr string) (err error) {
	var entry orm.Blacklist

	if entry = orm.GetBlacklistEntry(addr); entry.AddressId != 0 {
		err = errors.New("[Blacklist.Remove] Failed to retrieve address")
		return
	}

	bl.cache.Remove(addr)

	if ok := entry.Remove(); !ok {
		err = errors.New("[Blacklist.Remove]: Failed to remove entry")
		return
	}

	proto.RemoveBGPRoute(addr)

	return
}

func (bl *Blacklist) Listed(addr string) bool {
	return bl.cache.Has(addr)
}

func (bl *Blacklist) ReapExpiredEntries() {
	var cached_addr string

	for cached_addr, _ = range bl.cache.GetAll() {
		bl.Remove(cached_addr)
		Log.Debug(MYNAME + ": " + cached_addr + " expired from blacklist")
	}
}
