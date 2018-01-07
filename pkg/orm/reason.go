package orm

import (
	"fmt"

	"strings"
)

type Reason struct {
	Id     int64
	Reason string
}

func (obj *Reason) String() string {
	return fmt.Sprintf("Reason<%d %s>", obj.Id, obj.Reason)
}

func (o *ORM) GetReason(reason string) (*Reason, error) {
	entry := &Reason{}

	err := o.db.Model(entry).Where("reason = ?", reason).Select()
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("ORM.GetReason db.Select: %v", err)
		}
		err = nil
	}

	if entry.Reason == "" {
		return nil, nil
	}

	return entry, nil
}

func (o *ORM) GetReasonById(id int64) (*Reason, error) {
	entry := &Reason{}

	err := o.db.Model(entry).Where("id = ?", id).Select()
	if err != nil {
		return nil, fmt.Errorf("ORM.GetReasonById db.Select: %v", err)
	}

	return entry, nil
}

func (o *ORM) UpdateReason(reason_s string) (*Reason, error) {
	reason, err := o.GetReason(reason_s)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return nil, fmt.Errorf("ORM.UpdateReason: %v", err)
	}

	if reason == nil {
		o.log.Debugf("ORM.UpdateReason: adding new entry")
		reason = &Reason{
			Reason: reason_s,
		}

		err = o.db.Insert(reason)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateReason db.Create: %v", err)
		}

		q := Reason{Reason: reason_s}
		err = o.db.Select(q)
		if err != nil {
			return nil, fmt.Errorf("ORM.UpdateReason db.Select: %v", err)
		}

		reason = &q
	} else {
		o.log.Debugf("ORM.UpdateReason: updating existing entry")
		err = o.db.Update(reason)
		if err != nil {
			if !strings.Contains(err.Error(), "no rows in result set") {
				return nil, fmt.Errorf("ORM.UpdateReason db.Update: %v", err)
			}
			err = nil
		}
	}

	return reason, nil
}
