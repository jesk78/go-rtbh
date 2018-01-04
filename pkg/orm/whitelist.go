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
	err := db.Create(obj)
	if err != nil {
		return fmt.Errorf("Whitelist.Save db.Create: %v", err)
	}

	return nil
}

func (obj *Whitelist) Remove() error {
	err := db.Delete(obj)
	if err != nil {
		return fmt.Errorf("Whitelist.Remove db.Delete: %v", err)
	}

	return nil
}

func GetWhitelistEntry(addr_s string) (*Whitelist, error) {
	addr, err := GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry: %v", err)
	}
	if addr.Addr == "" {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry: address is emtry")
	}

	entry := &Whitelist{}
	err = db.Model(entry).Where(T_WHITELIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetWhitelistEntry db.Select: %v", err)
	}

	return entry, nil
}
