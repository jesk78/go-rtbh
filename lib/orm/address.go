package orm

import (
	"fmt"
)

type Address struct {
	Id   int64
	Addr string
	Fqdn string
}

func (obj Address) String() string {
	return fmt.Sprintf("Address<%d %s %s>", obj.Id, obj.Addr, obj.Fqdn)
}

func GetAddress(addr string) Address {
	var entry Address
	var err error

	err = Db.Model(&entry).Where("addr = ?", addr).Select()
	if err != nil {
		return Address{}
	}

	return entry
}

func GetAddressById(id int64) Address {
	var entry Address
	var err error

	err = Db.Model(&entry).Where("id = ?", id).Select()
	if err != nil {
		return Address{}
	}

	return entry
}

func GetAddressesNoFqdn() []Address {
	var addrs []Address
	var err error

	err = Db.Model(&addrs).Where("fqdn = ''").Select()
	if err != nil {
		return []Address{}
	}

	return addrs
}

func UpdateAddress(addr_s string, fqdn string) Address {
	var addr Address
	var err error

	addr = GetAddress(addr_s)
	if addr.Addr == "" {
		addr = Address{
			Addr: addr_s,
			Fqdn: fqdn,
		}

		err = Db.Create(&addr)
		if err != nil {
			Log.Warning("[orm]: UpdateAddress(" + addr_s + "," + fqdn + ") create failed: " + err.Error())
		}

	} else {
		err = Db.Update(&addr)
		if err != nil {
			Log.Warning("[orm]: UpdateAddress(" + addr_s + "," + fqdn + ") update failed: " + err.Error())
		}
	}

	return addr
}
