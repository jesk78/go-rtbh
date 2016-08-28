package lists

import (
	"errors"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/orm"
)

type History struct {
}

func (obj History) Add(event events.RTBHEvent) (err error) {
	var (
		addr   orm.Address
		entry  orm.History
		reason orm.Reason
	)

	if addr = orm.GetAddress(event.Address); addr.Addr == "" {
		err = errors.New("[History.Add]: GetAddress() failed: No such address")
		return
	}

	if reason = orm.GetReason(event.Reason); reason.Reason == "" {
		err = errors.New("[History.Add]: GetReason() failed: No such reason")
		return
	}

	entry = orm.History{
		Address: &addr,
		Reason:  &reason,
		AddedAt: event.AddedAt,
	}
	if ok := entry.Save(); !ok {
		err = errors.New("[History.Add]: " + entry.String() + ".Save() failed")
	}

	return
}
