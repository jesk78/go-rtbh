package orm

import (
	"fmt"
)

type Duration struct {
	Id       int64
	Duration string
}

func (obj Duration) String() string {
	return fmt.Sprintf("Duration<%d %s>", obj.Id, obj.Duration)
}

func GetDuration(duration string) Duration {
	var entry Duration
	var err error

	err = Db.Model(&entry).Where("duration = ?", duration).Select()
	if err != nil {
		return Duration{}
	}

	return entry
}

func UpdateDuration(duration_s string) Duration {
	var duration Duration
	var err error

	duration = GetDuration(duration_s)
	if duration.Duration == "" {
		duration = Duration{
			Duration: duration_s,
		}

		err = Db.Create(&duration)
		if err != nil {
			Log.Warning("[orm]: UpdateDuration(" + duration_s + ") create failed: " + err.Error())
		}
	} else {
		err = Db.Update(&duration)
		if err != nil {
			Log.Warning("[orm]: UpdateDuration(" + duration_s + ") update failed: " + err.Error())
		}
	}

	return duration
}
