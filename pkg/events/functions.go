package events

import (
	"encoding/json"
	"fmt"
	"time"
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
			etData := &etAlert{}
			err := json.Unmarshal(data, &etData)
			if err != nil {
				return fmt.Errorf("RTBHEvent.LoadFrom json.Unmarshal: %v", err)
			}
			event.Address = etData.SrcIp
			event.Reason = etData.Alert.Signature
			event.FlowId = etData.FlowId

			ts, err := time.Parse("2006-01-02T15:04:05.000000-0700", etData.Timestamp)
			if err != nil {
				return fmt.Errorf("RTBHEvent.LoadFrom time.Parse: %v", err)
			}
			event.AddedAt = ts
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
