package orm

import (
	"fmt"
	"time"
)

const BLACKLIST string = MYNAME + ".Blacklist"

type Blacklist struct {
	Id         int64
	AddrId     int64
	ReasonId   int64
	AddedAt    time.Time
	DurationId int64
}

func (obj *Blacklist) String() string {
	return fmt.Sprintf("Blacklist<%d %s %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt, obj.DurationId)
}

func (obj *Blacklist) Save() bool {
	var err error

	if err = db.Create(obj); err != nil {
		Log.Fatal(BLACKLIST + ": " + obj.String() + ".Save() failed: " + err.Error())
	}

	return true
}

func (obj *Blacklist) Remove() bool {
	var err error

	if err = db.Delete(obj); err != nil {
		Log.Fatal(BLACKLIST + ": " + obj.String() + ".Remove() failed: " + err.Error())
		return false
	}

	return true
}

func GetBlacklistEntry(addr_s string) *Blacklist {
	var addr *Address
	var entry *Blacklist
	var err error

	if addr = GetAddress(addr_s); addr == nil {
		Log.Warning(BLACKLIST + ": Address not found for " + addr_s)
		return nil
	}

	entry = &Blacklist{}
	err = db.Model(entry).Where(T_BLACKLIST+".addr_id = ?", addr.Id).Select()
	if err != nil {
		Log.Debug("err: " + err.Error() + " for addr " + addr.Addr)
		Log.Debug(entry)
		Log.Fatal(BLACKLIST + ": Failed to select blacklist entry for " + addr.String())
	}

	return entry
}

func GetBlacklistEntryByAddressId(addr_id int64) *Blacklist {
	var entry *Blacklist
	var err error

	entry = &Blacklist{}
	err = db.Model(entry).Where(T_BLACKLIST+".addr_id = ?", addr_id).Select()
	if err != nil {
		return nil
	}

	return entry
}
