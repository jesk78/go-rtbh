package pipeline

import (
	"encoding/json"
)

const E_ALERT string = "alert"
const E_HTTP string = "http"
const E_FILEINFO string = "fileinfo"
const E_SYSLOG string = "syslog"

// Struct used to extract the event type
type EventType struct {
	EventType string `json:"event_type"`
}

type etAlert struct {
	SrcIp string `json:"src_ip"`
	Alert struct {
		Signature string `json:"signature"`
	} `json:"alert"`
}

// The emitted event
type Event struct {
	Address string
	Reason  string
}

func (event *Event) LoadFrom(data []byte) (err error) {
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

func NewEvent(data []byte) (event *Event, err error) {
	event = &Event{}

	if err = event.LoadFrom(data); err != nil {
		return
	}

	return
}
