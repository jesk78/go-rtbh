package orm

import (
	"fmt"

	"strings"
)

type Address struct {
	Id   int64
	Addr string
	Fqdn string
}

func (obj Address) String() string {
	return fmt.Sprintf("Address<%d %s %s>", obj.Id, obj.Addr, obj.Fqdn)
}

func (o *ORM) GetAddress(addr string) (*Address, error) {
	entry := &Address{}

	err := o.db.Model(entry).Where("addr = ?", addr).Select()
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("ORM.GetAddress db.Select: %v", err)
		}
		err = nil
	}

	if entry.Addr == "" {
		return nil, nil
	}

	o.log.Debugf("ORM.GetAddress: got entry %v", entry)

	return entry, nil
}

func (o *ORM) GetAddressById(id int64) (*Address, error) {
	entry := &Address{}

	err := o.db.Model(entry).Where("id = ?", id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetAddressById db.Select: %v", err)
	}

	return entry, nil
}

func (o *ORM) GetAddressesNoFqdn() ([]Address, error) {
	addrs := []Address{}

	err := o.db.Model(&addrs).Where("fqdn = ''").Select()
	if err != nil {

		return nil, fmt.Errorf("ORM.GetAddressNoFqdn db.Select: %v", err)
	}

	return addrs, nil
}

func (o *ORM) UpdateAddress(addr_s, fqdn string) (*Address, error) {
	addr, err := o.GetAddress(addr_s)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return nil, fmt.Errorf("ORM.UpdateAddress: %v", err)
	}

	if addr == nil {
		o.log.Debugf("ORM.UpdateAddress: adding new entry")
		addr = &Address{
			Id:   0,
			Addr: addr_s,
			Fqdn: fqdn,
		}

		err = o.db.Insert(addr)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateAddress db.Create: %v", err)
		}

	} else {
		o.log.Debugf("ORM.UpdateAddress: updating existing entry")
		addr.Fqdn = fqdn
		err = o.db.Update(addr)
		if err != nil {
			if !strings.Contains(err.Error(), "no rows in result set") {
				return nil, fmt.Errorf("ORM.UpdateAddress db.Update: %v", err)
			}
			err = nil
		}
	}

	return addr, nil
}
