package orm

import (
	"fmt"
)

type Whitelist struct {
	Id          int64
	Address     *Address
	Description string
}

func (obj Whitelist) String() string {
	return fmt.Sprintf("Whitelist<%d %s %s>", obj.Id, obj.Address.Addr, obj.Description)
}

func (obj Whitelist) Save() bool {
	var err error

	if err = db.Create(&obj); err != nil {
		Log.Warning("[orm] " + obj.String() + ".Save() failed: " + err.Error())
		return false
	}

	return true
}

func (obj Whitelist) Remove() bool {
	var err error

	if err = db.Delete(&obj); err != nil {
		Log.Warning("[orm] " + obj.String() + ".Remove() failed: " + err.Error())
		return false
	}
	return true
}

func GetWhitelistEntry(addr_s string) Whitelist {
	var addr Address
	var entry Whitelist
	var err error

	if addr = GetAddress(addr_s); addr.Addr == "" {
		return Whitelist{}
	}

	err = db.Model(&entry).Where(T_WHITELIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		return Whitelist{}
	}

	return entry
}
