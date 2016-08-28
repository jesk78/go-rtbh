package orm

import (
	"fmt"
	"time"
)

type Blacklist struct {
	Id         int64
	AddressId  int64
	ReasonId   int64
	AddedAt    time.Time
	DurationId int64
}

func (obj Blacklist) String() string {
	return fmt.Sprintf("Blacklist<%d %s %s %s %s>", obj.Id, obj.AddressId, obj.ReasonId, obj.AddedAt, obj.DurationId)
}

func (obj Blacklist) Save() bool {
	var err error

	if err = Db.Create(&obj); err != nil {
		Log.Warning("[orm] " + obj.String() + ".Save() failed: " + err.Error())
		return false
	}

	return true
}

func (obj Blacklist) Remove() bool {
	var err error

	if err = Db.Delete(&obj); err != nil {
		Log.Warning("[orm] " + obj.String() + ".Remove() failed: " + err.Error())
		return false
	}

	return true
}

func GetBlacklistEntry(addr_s string) Blacklist {
	var addr Address
	var entry Blacklist
	var err error

	if addr = GetAddress(addr_s); addr.Addr == "" {
		return Blacklist{}
	}

	err = Db.Model(&entry).Where(T_BLACKLIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		return Blacklist{}
	}

	return entry
}
