package orm

import (
	"fmt"
	"time"

	"github.com/r3boot/go-rtbh/pkg/memcache"
)

type History struct {
	Id       int64
	AddrId   int64
	ReasonId int64
	AddedAt  time.Time
}

var cacheHistory *memcache.StringCache

func (obj *History) String() string {
	return fmt.Sprintf("History<%d %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt)
}

func (obj *History) Save() error {
	addr, err := GetAddressById(obj.AddrId)
	if err != nil {
		return fmt.Errorf("History.Save: %v", err)
	}
	if addr.Addr == "" {
		return fmt.Errorf("History.Save: No address record found for %v", obj)
	}

	err = db.Create(obj)
	if err != nil {
		return fmt.Errorf("History.Save db.Create: %v", err)
	}

	cacheHistory.Add(addr.Addr, obj)

	return nil
}

func GetHistoryEntries(addr_s string) ([]*History, error) {
	addr, err := GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.GetHistoryEntries: %v", err)
	}
	if addr.Addr == "" {
		return nil, fmt.Errorf("ORM.GetHistoryEntries: No such address")
	}

	entries := []*History{}
	err = db.Model(entries).Where("?.addr_id = ?", T_HISTORY, addr.Id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetHistoryEntries db.Select: %v", err)
	}

	return entries, nil
}
