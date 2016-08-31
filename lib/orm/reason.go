package orm

import (
	"fmt"
)

const REASON string = MYNAME + ".Reason"

type Reason struct {
	Id     int64
	Reason string
}

func (obj *Reason) String() string {
	return fmt.Sprintf("Reason<%d %s>", obj.Id, obj.Reason)
}

func GetReason(reason string) *Reason {
	var entry *Reason
	var err error

	entry = &Reason{}
	err = db.Model(entry).Where("reason = ?", reason).Select()
	if err != nil {
		Log.Debug(REASON + ": Failed to lookup reason for " + reason + ": " + err.Error())
		return nil
	}

	return entry
}

func UpdateReason(reason_s string) *Reason {
	var reason *Reason
	var err error

	reason = GetReason(reason_s)
	if reason == nil {
		reason = &Reason{
			Reason: reason_s,
		}

		err = db.Create(reason)
		if err != nil {
			Log.Fatal(REASON + ".UpdateReason(" + reason_s + ") create failed: " + err.Error())
		}
	} else {
		err = db.Update(reason)
		if err != nil {
			Log.Warning(REASON + ".UpdateReason(" + reason_s + ") update failed: " + err.Error())
		}
	}

	return reason
}
