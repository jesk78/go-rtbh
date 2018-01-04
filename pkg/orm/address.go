package orm

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/memcache"
)

type Address struct {
	Id   int64
	Addr string
	Fqdn string
}

var (
	cacheAddrOnAddress *memcache.StringCache
	cacheAddrOnId      *memcache.IntCache
)

func (obj Address) String() string {
	return fmt.Sprintf("Address<%d %s %s>", obj.Id, obj.Addr, obj.Fqdn)
}

func GetAddress(addr string) (*Address, error) {
	entry := &Address{}

	if cacheAddrOnAddress.Has(addr) {
		tmp := cacheAddrOnAddress.Get(addr).(Address)
		entry = &tmp
	} else {
		err := db.Model(entry).Where("addr = ?", addr).Select()
		if err != nil {
			return nil, fmt.Errorf("ORM.GetAddress db.Select: %v", err)
		}
		cacheAddrOnAddress.Add(addr, entry)
	}

	return entry, nil
}

func GetAddressById(id int64) (*Address, error) {
	entry := &Address{}

	if cacheAddrOnId.Has(id) {
		tmp := cacheAddrOnId.Get(id).(Address)
		entry = &tmp
	} else {
		err := db.Model(entry).Where("id = ?", id).Select()
		if err != nil {
			return nil, fmt.Errorf("ORM.GetAddressById db.Select: %v", err)
		}
		cacheAddrOnId.Add(id, entry)
	}

	return entry, nil
}

func GetAddressesNoFqdn() ([]Address, error) {
	addrs := []Address{}

	err := db.Model(&addrs).Where("fqdn = ''").Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetAddressNoFqdn db.Select: %v", err)
	}

	return addrs, nil
}

func UpdateAddress(addr_s string, fqdn string) (*Address, error) {
	addr, err := GetAddress(addr_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.UpdateAddress: %v", err)
	}

	if addr.Addr == "" {
		addr = &Address{
			Id:   0,
			Addr: addr_s,
			Fqdn: fqdn,
		}

		err = db.Create(addr)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateAddress db.Create: %v", err)
		}
		cacheAddrOnAddress.Add(addr_s, addr)
		// TODO: Add cacheAddrOnId

	} else {
		addr.Fqdn = fqdn
		err = db.Update(addr)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateAddress db.Update: %v", err)
		}
		cacheAddrOnAddress.Add(addr_s, addr)
		cacheAddrOnId.Add(addr.Id, addr)
	}

	return addr, nil
}

func WarmupAddressCaches() error {
	addresses := []Address{}
	_, err := db.Query(&addresses, "SELECT * FROM addresses")
	if err != nil {
		return fmt.Errorf("ORM.WarmupAddressCaches db.Query: %v", err)
	}

	for _, addr := range addresses {
		cacheAddrOnAddress.Add(addr.Addr, addr)
		cacheAddrOnId.Add(addr.Id, addr)
	}

	return nil
}
