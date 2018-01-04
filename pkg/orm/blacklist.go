package orm

import (
	"fmt"
	"time"

	"github.com/r3boot/go-rtbh/pkg/memcache"
)

type Blacklist struct {
	Id       int64
	AddrId   int64
	ReasonId int64
	AddedAt  time.Time
	ExpireOn time.Time
}

var (
	cacheBlacklistOnAddress *memcache.StringCache
	cacheBlacklistOnId      *memcache.IntCache
)

func (obj *Blacklist) String() string {
	return fmt.Sprintf("Blacklist<%d %s %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt, obj.ExpireOn)
}

func (obj *Blacklist) Save(addr_s string) error {
	err := db.Create(obj)
	if err != nil {
		return fmt.Errorf("Blacklist.Save db.Create: %v", err)
	}

	cacheBlacklistOnAddress.Add(addr_s, obj)
	cacheBlacklistOnId.Add(obj.Id, obj)

	return nil
}

func (obj *Blacklist) Remove() error {
	addr, err := GetAddressById(obj.AddrId)
	if err != nil {
		return fmt.Errorf("Blacklist.Remove: %v", err)
	}

	if addr.Addr == "" {
		return fmt.Errorf("Blacklist.Remove: No address record found for object id %d", obj.Id)
	}

	cacheBlacklistOnAddress.Remove(addr.Addr)
	cacheBlacklistOnId.Remove(obj.Id)

	if err = db.Delete(obj); err != nil {
		return fmt.Errorf("Blacklist.Remove db.Delete: %v", err)
	}

	return nil
}

func GetBlacklistEntry(addr_s string) (*Blacklist, error) {
	entry := &Blacklist{}

	if cacheBlacklistOnAddress.Has(addr_s) {
		entry = cacheBlacklistOnAddress.Get(addr_s).(*Blacklist)
	} else {
		addr, err := GetAddress(addr_s)
		if addr.Addr == "" {
			return nil, fmt.Errorf("ORM.GetBlacklistEntry: %v", err)
		}

		err = db.Model(entry).Where(T_BLACKLIST+".addr_id = ?", addr.Id).Select()
		if err != nil {
			return nil, fmt.Errorf("ORM.GetBlacklistEntry db.Select: %v", err)
		}

		cacheBlacklistOnAddress.Add(addr_s, entry)
	}

	return entry, nil
}

func GetBlacklistEntries() []Blacklist {
	entries := []Blacklist{}

	for _, entry := range cacheBlacklistOnAddress.GetAll() {
		entries = append(entries, entry.(Blacklist))
	}

	return entries
}

func WarmupBlacklistCaches() error {
	entries := []Blacklist{}
	_, err := db.Query(&entries, "SELECT * FROM blacklists")
	if err != nil {
		return fmt.Errorf("ORM.WarmupBlacklistCaches db.Query: %v", err)
	}

	for _, entry := range entries {
		addr, err := GetAddressById(entry.AddrId)
		if err != nil {
			return fmt.Errorf("ORM.WarmupBlacklistCaches: %v", err)
		}

		if addr.Addr == "" {
			return fmt.Errorf("ORM.WarmupBlacklistCaches: Failed to retrieve address for addr_id %", entry.AddrId)
		}

		cacheBlacklistOnAddress.Add(addr.Addr, entry)
		cacheBlacklistOnId.Add(entry.AddrId, entry)
	}

	return nil
}
