package events

import (
	"time"
)

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

// Struct containing an event which is being processed
type RTBHEvent struct {
	Address  string
	Reason   string
	AddedAt  time.Time
	ExpireIn string
}

// Struct containing a whitelist entry
type RTBHWhiteEntry struct {
	Address     string
	Description string
}
