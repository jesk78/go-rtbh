package events

import (
	"encoding/json"
	"time"
)

func (event *RTBHEvent) LoadFrom(data []byte) (err error) {
	var et EventType

	// First, determine the event type
	if err = json.Unmarshal(data, &et); err != nil {
		return
	}

	switch et.EventType {
	case E_ALERT:
		{
			et_data := &etAlert{}
			if err = json.Unmarshal(data, &et_data); err != nil {
				return
			}
			event.Address = et_data.SrcIp
			event.Reason = et_data.Alert.Signature
		}
	case E_HTTP:
		{
		}
	case E_FILEINFO:
		{
		}
	case E_SYSLOG:
		{
		}
	default:
		{
			Log.Debug("Unknown EventType: " + et.EventType)
		}
	}
	return
}

func NewRTBHEvent(data []byte) (event *RTBHEvent, err error) {
	event = &RTBHEvent{
		AddedAt: time.Now(),
	}

	if err = event.LoadFrom(data); err != nil {
		return
	}

	return
}
