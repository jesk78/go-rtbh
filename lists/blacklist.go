package lists

import (
	"errors"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/orm"
	"github.com/r3boot/go-rtbh/proto"
	"net"
	"time"
)

type cacheBlacklist struct {
	Reason   string
	AddedAt  time.Time
	Duration string
}

type Blacklist struct {
	Cache map[string]*cacheBlacklist
}

func NewBlacklist() *Blacklist {
	var blacklist *Blacklist

	blacklist = &Blacklist{}
	blacklist.Cache = make(map[string]*cacheBlacklist)

	return blacklist
}

func (obj Blacklist) cacheHas(addr string) bool {
	var cached_addr string

	for cached_addr, _ = range obj.Cache {
		if cached_addr == addr {
			return true
		}
	}

	return false
}

func (obj Blacklist) cacheUpdate(event events.RTBHEvent) *cacheBlacklist {
	var cached_entry *cacheBlacklist

	cached_entry = &cacheBlacklist{
		Reason:   event.Reason,
		AddedAt:  event.AddedAt,
		Duration: event.ExpireIn,
	}

	obj.Cache[event.Address] = cached_entry

	return cached_entry
}

func (obj Blacklist) cachePurge(addr string) {
	if ok := obj.cacheHas(addr); !ok {
		return
	}

	delete(obj.Cache, addr)
}

func (obj Blacklist) Add(event events.RTBHEvent) (err error) {
	var (
		addr     orm.Address
		duration orm.Duration
		entry    orm.Blacklist
		names    []string
		fqdn     string
		reason   orm.Reason
	)

	if names, err = net.LookupAddr(event.Address); err != nil {
		Log.Warning("[Blacklist]: Failed to lookup fqdn for " + event.Address)
		fqdn = ""
	} else {
		fqdn = names[0]
	}

	if len(names) > 1 {
		Log.Warning("[Blacklist]: Multiple hostnames found for " + event.Address + ", using: " + fqdn)
	}

	if addr = orm.UpdateAddress(event.Address, fqdn); addr.Addr == "" {
		Log.Debug("[Blacklist]: Failed to update Address record")
		return
	}

	if duration = orm.UpdateDuration(event.ExpireIn); duration.Duration == "" {
		Log.Debug("[Blacklist]: Failed to update Duration record")
		return
	}

	if reason = orm.UpdateReason(event.Reason); reason.Reason == "" {
		Log.Debug("[Blacklist]: Failed to update Reason record")
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

	proto.AddBGPRoute(addr.Addr)

	obj.cacheUpdate(event)

	return
}

func (obj Blacklist) Remove(addr string) (err error) {
	var entry orm.Blacklist

	if entry = orm.GetBlacklistEntry(addr); entry.AddressId != 0 {
		err = errors.New("[Blacklist.Remove] Failed to retrieve address")
		return
	}

	if ok := entry.Remove(); !ok {
		err = errors.New("[Blacklist.Remove]: Failed to remove entry")
		return
	}

	proto.RemoveBGPRoute(addr)

	return
}

func (obj Blacklist) Listed(addr string) bool {
	return obj.cacheHas(addr)
}
