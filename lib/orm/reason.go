package orm

import (
	"fmt"

	"github.com/r3boot/go-rtbh/lib/memcache"
)

type Reason struct {
	Id     int64
	Reason string
}

var (
	cacheReasonOnAddress *memcache.StringCache
	cacheReasonOnId      *memcache.IntCache
)

func (obj *Reason) String() string {
	return fmt.Sprintf("Reason<%d %s>", obj.Id, obj.Reason)
}

func GetReason(reason string) (*Reason, error) {
	entry := &Reason{}
	if cacheReasonOnAddress.Has(reason) {
		tmp := cacheReasonOnAddress.Get(reason).(Reason)
		entry = &tmp
	} else {
		err := db.Model(entry).Where("reason = ?", reason).Select()
		if err != nil {
			return nil, fmt.Errorf("ORM.GetReason db.Select: %v", err)
		}

		cacheReasonOnAddress.Add(reason, entry)
	}

	return entry, nil
}

func GetReasonById(id int64) (*Reason, error) {
	entry := &Reason{}

	if cacheReasonOnId.Has(id) {
		tmp := cacheReasonOnId.Get(id).(Reason)
		entry = &tmp
	} else {
		err := db.Model(entry).Where("id = ?", id).Select()
		if err != nil {
			return nil, fmt.Errorf("ORM.GetReasonById db.Select: %v", err)
		}
		cacheReasonOnId.Add(id, entry)
	}

	return entry, nil
}

func UpdateReason(reason_s string) (*Reason, error) {
	reason, err := GetReason(reason_s)
	if err != nil {
		return nil, fmt.Errorf("ORM.UpdateReason: %v", err)
	}

	if reason.Reason == "" {
		reason = &Reason{
			Reason: reason_s,
		}

		err = db.Create(reason)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateReason db.Create: %v", err)
		}
	} else {
		err = db.Update(reason)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateReason db.Update: %v", err)
		}
	}

	cacheReasonOnAddress.Add(reason_s, reason)
	cacheReasonOnId.Add(reason.Id, reason)

	return reason, nil
}

func WarmupReasonCaches() error {
	reasons := []Reason{}
	_, err := db.Query(&reasons, "SELECT * FROM reasons")
	if err != nil {
		return fmt.Errorf("ORM.WarmupReasonCaches db.Query: %v", err)
	}

	for _, reason := range reasons {
		cacheReasonOnAddress.Add(reason.Reason, reason)
		cacheReasonOnId.Add(reason.Id, reason)
	}

	return nil
}
