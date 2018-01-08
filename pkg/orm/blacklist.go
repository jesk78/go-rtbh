package orm

import (
	"fmt"
	"strings"
	"time"
)

type Blacklist struct {
	Id       int64
	AddrId   int64
	ReasonId int64
	FlowId   int64
	AddedAt  time.Time
	ExpireOn time.Time
}

func (obj *Blacklist) String() string {
	return fmt.Sprintf("Blacklist<%d %s %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt, obj.ExpireOn)
}

func (obj *Blacklist) Save(addr_s string) error {
	err := localORM.db.Insert(obj)
	if err != nil {
		return fmt.Errorf("Blacklist.Save db.Create: %v", err)
	}

	return nil
}

func (obj *Blacklist) Remove() error {
	addr, err := localORM.GetAddressById(obj.AddrId)
	if err != nil {
		return fmt.Errorf("Blacklist.Remove: %v", err)
	}

	if addr.Addr == "" {
		return fmt.Errorf("Blacklist.Remove: No address record found for object id %d", obj.Id)
	}

	err = localORM.db.Delete(obj)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return fmt.Errorf("Blacklist.Remove db.Delete: %v", err)
	}

	return nil
}

func (o *ORM) GetBlacklistEntry(addr_s string) (*Blacklist, error) {
	entry := &Blacklist{}

	addr, err := o.GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.GetBlacklistEntry: %v", err)
	}

	if (addr == nil) || (addr != nil && addr.Addr == "") {
		return nil, fmt.Errorf("ORM.GetBlacklistEntry: address is empty")
	}

	err = o.db.Model(entry).Where(T_BLACKLIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetBlacklistEntry db.Select: %v", err)
	}

	return entry, nil
}

func (o *ORM) GetBlacklistEntries() ([]*Blacklist, error) {
	entries := []*Blacklist{}

	_, err := o.db.Query(&entries, "SELECT * FROM blacklists ORDER BY added_at DESC")
	if err != nil {
		return nil, fmt.Errorf("ORM.GetBlacklistEntries db.Query: %v", err)
	}

	return entries, nil
}
