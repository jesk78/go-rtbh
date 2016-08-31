package history

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/orm"
)

func (history *History) Add(event events.RTBHEvent) (err error) {
	var (
		addr   *orm.Address
		entry  *orm.History
		reason *orm.Reason
	)

	if addr = orm.GetAddress(event.Address); addr == nil {
		err = errors.New(MYNAME + ": GetAddress() failed: No such address")
		return
	}

	if reason = orm.GetReason(event.Reason); reason.Reason == "" {
		err = errors.New(MYNAME + ": GetReason() failed: No such reason")
		return
	}

	Log.Debug("Working on " + event.Address)
	entry = &orm.History{
		AddrId: addr.Id,
		ReasonId:    reason.Id,
		AddedAt:     event.AddedAt,
	}
	if ok := entry.Save(); !ok {
		err = errors.New(MYNAME + ": " + entry.String() + ".Save() failed")
	}

	return
}
