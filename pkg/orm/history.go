package orm

import (
	"fmt"
	"time"
)

type History struct {
	Id       int64
	AddrId   int64
	ReasonId int64
	AddedAt  time.Time
}

func (obj *History) String() string {
	return fmt.Sprintf("History<%d %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt)
}

func (obj *History) Save() error {
	addr, err := localORM.GetAddressById(obj.AddrId)
	if err != nil {
		return fmt.Errorf("History.Save: %v", err)
	}
	if addr.Addr == "" {
		return fmt.Errorf("History.Save: No address record found for %v", obj)
	}

	err = localORM.db.Insert(obj)
	if err != nil {
		return fmt.Errorf("History.Save db.Create: %v", err)
	}

	return nil
}

func (o *ORM) GetHistoryEntries(addr_s string) ([]*History, error) {
	addr, err := o.GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.GetHistoryEntries: %v", err)
	}
	if addr.Addr == "" {
		return nil, fmt.Errorf("ORM.GetHistoryEntries: No such address")
	}

	entries := []*History{}
	err = o.db.Model(entries).Where("?.addr_id = ?", T_HISTORY, addr.Id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetHistoryEntries db.Select: %v", err)
	}

	return entries, nil
}
