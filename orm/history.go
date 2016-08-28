package orm

import (
	"fmt"
	"time"
)

type History struct {
	Id      int64
	Address *Address
	Reason  *Reason
	AddedAt time.Time
}

func (obj History) String() string {
	return fmt.Sprintf("History<%d %s %s %s>", obj.Id, obj.Address.Addr, obj.Reason.Reason, obj.AddedAt)
}

func (obj History) Save() bool {
	var err error

	if err = Db.Create(&obj); err != nil {
		Log.Warning("[orm]: " + obj.String() + ".Save() failed: " + err.Error())
		return false
	}

	return true
}

func GetHistoryEntries(addr_s string) []History {
	var addr Address
	var entries []History
	var err error

	if addr = GetAddress(addr_s); addr.Addr == "" {
		Log.Warning("[orm]: GetHistoryEntries(" + addr_s + ") failed: GetAddress(): No such address")
		return []History{}
	}

	err = Db.Model(&entries).Where("?.addr_id = ?", T_HISTORY, addr.Id).Select()
	if err != nil {
		Log.Warning("[orm]: GetHistoryEntries(" + addr_s + ") failed: " + err.Error())
		return []History{}
	}

	return entries
}
