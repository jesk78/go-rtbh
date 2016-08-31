package orm

import (
	"fmt"
	"strconv"
)

const ADDRESS string = MYNAME + ".Address"

type Address struct {
	Id   int64
	Addr string
	Fqdn string
}

func (obj *Address) String() string {
	return fmt.Sprintf("Address<%d %s %s>", obj.Id, obj.Addr, obj.Fqdn)
}

func GetAddress(addr string) *Address {
	var entry *Address
	var err error

	entry = &Address{}

	err = db.Model(entry).Where("addr = ?", addr).Select()
	if err != nil {
		Log.Debug(ADDRESS + ": Failed to lookup address for " + addr + ": " + err.Error())
		return nil
	}

	return entry
}

func GetAddressById(id int64) *Address {
	var entry *Address
	var err error

	entry = &Address{}

	err = db.Model(entry).Where("id = ?", id).Select()
	if err != nil {
		Log.Debug(ADDRESS + ": Failed to lookup address for " + strconv.Itoa(int(id)))
		return nil
	}

	return entry
}

func GetAddressesNoFqdn() []*Address {
	var addrs []*Address
	var err error

	err = db.Model(addrs).Where("fqdn = ''").Select()
	if err != nil {
		return nil
	}

	return addrs
}

func UpdateAddress(addr_s string, fqdn string) *Address {
	var addr *Address
	var err error

	addr = GetAddress(addr_s)
	if addr == nil {
		addr = &Address{
			Id:   0,
			Addr: addr_s,
			Fqdn: fqdn,
		}

		err = db.Create(addr)
		if err != nil {
			Log.Debug(addr.String())
			Log.Fatal(ADDRESS + ".UpdateAddress(" + addr_s + "," + fqdn + ") create failed: " + err.Error())
		}

	} else {
		addr.Fqdn = fqdn
		err = db.Update(addr)
		if err != nil {
			Log.Fatal(ADDRESS + ".UpdateAddress(" + addr_s + "," + fqdn + ") update failed: " + err.Error())
		}
	}

	return addr
}
