package es

import (
	"net/http"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

type AlertGeoIpLocationData struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type AlertGeoIpData struct {
	ContinentCode string                 `json:"continent_code"`
	CountryCode   string                 `json:"country_code3"`
	CountryName   string                 `json:"country_name"`
	Ip            string                 `json:"ip"`
	Latitude      float32                `json:"latitude"`
	Location      AlertGeoIpLocationData `json:"location"`
	Longitude     float32                `json:"longitude"`
	Timezone      string                 `json:"timezone"`
}

type AlertMetaData struct {
	Gid         int32  `json:"gid"`
	Severity    int32  `json:"severity"`
	Signature   string `json:"signature"`
	Action      string `json:"action"`
	Category    string `json:"category"`
	SignatureId int32  `json:"signature_id"`
	Rev         int32  `json:"rev"`
}

type AlertFlowData struct {
	BytesToServer int64  `json:"bytes_toserver"`
	BytesToClient int64  `json:"bytes_toclient"`
	PktsToServer  int64  `json:"pkts_toserver"`
	PktsToClient  int64  `json:"pkts_toclient"`
	Start         string `json:"start"`
}

type AlertData struct {
	AppProto  string         `json:"app_proto"`
	DstIp     string         `json:"dest_ip"`
	DstPort   int32          `json:"dest_port"`
	DstGeoIp  AlertGeoIpData `json:"geoip_dstip"`
	FlowId    int64          `json:"flow_id"`
	Proto     string         `json:"proto"`
	Timestamp string         `json:"timestamp"`
	Sensor    string         `json:"host"`
	SrcIp     string         `json:"src_ip"`
	SrcPort   int32          `json:"src_port"`
	SrcGeoIp  AlertGeoIpData `json:"geoip_src_ip"`
	Alert     AlertMetaData  `json:"alert"`
	Flow      AlertFlowData  `json:"flow"`
}

type ES struct {
	cfg    *config.Config
	log    *logger.Logger
	client *http.Client
}
