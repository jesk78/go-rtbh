package events

import "time"

const (
	E_ALERT    = "alert"
	E_HTTP     = "http"
	E_FILEINFO = "fileinfo"
	E_SYSLOG   = "syslog"
)

type eventTypeDiscovery struct {
	EventType string `json:"event_type"`
}

type etAlert struct {
	SrcIp string `json:"src_ip"`
	FlowId int64 `json:"flow_id"`
	Timestamp string `json:"timestamp"`
	Alert struct {
		Signature string `json:"signature"`
	} `json:"alert"`
}

type RTBHEvent struct {
	Address  string
	Reason   string
	FlowId 	 int64
	AddedAt  time.Time
	ExpireIn string
}

type APIEvent struct {
	Id       int64
	Address  string
	Reason   string
	FlowId int64
	AddedAt  time.Time
	ExpireOn time.Time
}

// Struct containing a whitelist entry
type RTBHWhiteEntry struct {
	Id int64
	Address     string
	Description string
}
