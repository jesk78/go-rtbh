package orm

import (
	"fmt"
)

type Whitelist struct {
	Id          int64
	AddrId      int64
	Description string
}

func (obj *Whitelist) String() string {
	return fmt.Sprintf("Whitelist<%d %s %s>", obj.Id, obj.AddrId, obj.Description)
}

func (obj *Whitelist) Save() error {
	err := localORM.db.Insert(obj)
	if err != nil {
		return fmt.Errorf("Whitelist.Save db.Create: %v", err)
	}

	return nil
}

func (obj *Whitelist) Update() error {
	err := localORM.db.Update(obj)
	if err != nil {
		return fmt.Errorf("Whitelist.Update db.Create: %v", err)
	}

	return nil
}

func (obj *Whitelist) Remove() error {
	err := localORM.db.Delete(obj)
	if err != nil {
		return fmt.Errorf("Whitelist.Remove db.Delete: %v", err)
	}

	return nil
}

func (o *ORM) GetWhitelistEntry(addr_s string) (*Whitelist, error) {
	addr, err := o.GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry: %v", err)
	}

	if (addr == nil) || (addr != nil && addr.Addr == "") {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry: address is emtry")
	}

	entry := &Whitelist{}
	err = o.db.Model(entry).Where(T_WHITELIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry db.Select: %v", err)
	}

	return entry, nil
}

func (o *ORM) GetWhitelistEntries() ([]*Whitelist, error) {
	entries := []*Whitelist{}

	_, err := o.db.Query(&entries, "SELECT * FROM whitelists")
	if err != nil {
		return nil, fmt.Errorf("ORM.GetWhitelistEntries db.Query: %v", err)
	}

	return entries, nil
}
