package orm

import (
	"fmt"
)

type Reason struct {
	Id     int64
	Reason string
}

func (obj Reason) String() string {
	return fmt.Sprintf("Reason<%d %s>", obj.Id, obj.Reason)
}

func GetReason(reason string) Reason {
	var entry Reason
	var err error

	err = db.Model(&entry).Where("reason = ?", reason).Select()
	if err != nil {
		return Reason{}
	}

	return entry
}

func UpdateReason(reason_s string) Reason {
	var reason Reason
	var err error

	reason = GetReason(reason_s)
	if reason.Reason == "" {
		reason = Reason{
			Reason: reason_s,
		}

		err = db.Create(&reason)
		if err != nil {
			Log.Warning("[orm]: UpdateReason(" + reason_s + ") create failed: " + err.Error())
		}
	} else {
		err = db.Update(&reason)
		if err != nil {
			Log.Warning("[orm]: UpdateReason(" + reason_s + ") update failed: " + err.Error())
		}
	}

	return reason
}
