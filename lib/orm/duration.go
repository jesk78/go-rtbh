package orm

import (
	"fmt"
)

const DURATION string = MYNAME + ".Duration"

type Duration struct {
	Id       int64
	Duration string
}

func (obj *Duration) String() string {
	return fmt.Sprintf("Duration<%d %s>", obj.Id, obj.Duration)
}

func GetDuration(duration string) *Duration {
	var entry *Duration
	var err error

	entry = &Duration{}
	err = db.Model(entry).Where("duration = ?", duration).Select()
	if err != nil {
		Log.Debug(DURATION + ": Failed to lookup duration for " + duration + ": " + err.Error())
		return nil
	}

	return entry
}

func UpdateDuration(duration_s string) *Duration {
	var duration *Duration
	var err error

	duration = GetDuration(duration_s)
	if duration == nil {
		duration = &Duration{
			Duration: duration_s,
		}

		err = db.Create(duration)
		if err != nil {
			Log.Fatal(DURATION + ".UpdateDuration(" + duration_s + ") create failed: " + err.Error())
		}
	} else {
		err = db.Update(duration)
		if err != nil {
			Log.Fatal(DURATION + ".UpdateDuration(" + duration_s + ") update failed: " + err.Error())
		}
	}

	return duration
}
