package orm

import (
	"fmt"
	"time"
)

const HISTORY string = MYNAME + ".History"

type History struct {
	Id       int64
	AddrId   int64
	ReasonId int64
	AddedAt  time.Time
}

func (obj *History) String() string {
	return fmt.Sprintf("History<%d %s %s %s>", obj.Id, obj.AddrId, obj.ReasonId, obj.AddedAt)
}

func (obj *History) Save() bool {
	var err error

	if err = db.Create(obj); err != nil {
		Log.Fatal(MYNAME + ": " + obj.String() + ".Save() failed: " + err.Error())
	}

	return true
}

func GetHistoryEntries(addr_s string) []*History {
	var addr *Address
	var entries []*History
	var err error

	if addr = GetAddress(addr_s); addr.Addr == "" {
		Log.Fatal(MYNAME + ".GetHistoryEntries(" + addr_s + ") failed: GetAddress(): No such address")
	}

	err = db.Model(entries).Where("?.addr_id = ?", T_HISTORY, addr.Id).Select()
	if err != nil {
		Log.Fatal(MYNAME + ".GetHistoryEntries(" + addr_s + ") failed: " + err.Error())
	}

	return entries
}
