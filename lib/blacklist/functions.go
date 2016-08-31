package blacklist

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/resolver"
)

const SUBNAME string = MYNAME + ".Blacklist"

func (bl *Blacklist) Add(event events.RTBHEvent) (err error) {
	var (
		addr     *orm.Address
		duration *orm.Duration
		entry    *orm.Blacklist
		fqdn     string
		reason   *orm.Reason
	)

	bl.mutex.Lock()

	// Set fqdn to not-yet-lookedup so it will be picked up by the Resolver
	fqdn = resolver.FQDN_TO_LOOKUP

	if addr = orm.UpdateAddress(event.Address, fqdn); addr.Addr == "" {
		err = errors.New(MYNAME + ": Failed to update Address record")
		bl.mutex.Unlock()
		return
	}

	if duration = orm.UpdateDuration(event.ExpireIn); duration.Duration == "" {
		err = errors.New(SUBNAME + ": Failed to update Duration record")
		bl.mutex.Unlock()
		return
	}

	if reason = orm.UpdateReason(event.Reason); reason.Reason == "" {
		err = errors.New(SUBNAME + ": Failed to update Reason record")
		bl.mutex.Unlock()
		return
	}

	entry = orm.GetBlacklistEntryByAddressId(addr.Id)
	if entry != nil {
		Log.Warning(SUBNAME + ": Entry for " + event.Address + " / " + event.Reason + " already exists")
		bl.mutex.Unlock()
		return
	}

	entry = &orm.Blacklist{
		AddrId:     addr.Id,
		ReasonId:   reason.Id,
		AddedAt:    event.AddedAt,
		DurationId: duration.Id,
	}
	if ok := entry.Save(); !ok {
		bl.mutex.Unlock()
		return
	}

	Log.Debug(SUBNAME + ": Adding BGP route")
	bl.bgp.AddRoute(addr.Addr)

	bl.cache.Add(addr.Addr, entry)

	bl.mutex.Unlock()

	return
}

func (bl *Blacklist) Remove(addr string) (err error) {
	var entry *orm.Blacklist

	if entry = orm.GetBlacklistEntry(addr); entry == nil {
		err = errors.New(SUBNAME + " Failed to locate " + addr + " on the blacklist")
		return
	}

	bl.cache.Remove(addr)

	if ok := entry.Remove(); !ok {
		err = errors.New(SUBNAME + ": Failed to remove " + addr + " from the blacklist")
		return
	}

	bl.bgp.RemoveRoute(addr)

	return
}

func (bl *Blacklist) Listed(addr string) bool {
	return bl.cache.Has(addr)
}

func (bl *Blacklist) ReapExpiredEntries() {
	var cached_addr string
	var err error

	for cached_addr, _ = range bl.cache.GetAll() {
		if err = bl.Remove(cached_addr); err != nil {
			Log.Warning(err)
			continue
		}
		Log.Debug(SUBNAME + ": " + cached_addr + " expired from blacklist")
	}
}
