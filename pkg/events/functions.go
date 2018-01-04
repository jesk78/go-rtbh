package events

import (
	"encoding/json"
	"fmt"
)

func (event *RTBHEvent) LoadFrom(data []byte) error {
	var et eventTypeDiscovery

	// First, determine the event type
	err := json.Unmarshal(data, &et)
	if err != nil {
		return fmt.Errorf("RTBHEvent.LoadFrom json.Unmarshal: %v", err)
	}

	switch et.EventType {
	case E_ALERT:
		{
			et_data := &etAlert{}
			err := json.Unmarshal(data, &et_data)
			if err != nil {
				return fmt.Errorf("RTBHEvent.LoadFrom json.Unmarshal: %v", err)
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
			log.Warningf("RTBHEvent.LoadFrom: Unknown EventType: %s", et.EventType)
		}
	}

	return nil
}
