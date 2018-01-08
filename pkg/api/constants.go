package api

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/es"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
)

const (
	PATH_API_BLACKLIST  = "/v1/blacklist"
	PATH_APP            = "/app/"
	PATH_LIBS           = "/node_modules/"
	PATH_VIEW_BLACKLIST = "/blacklist"
	PATH_VIEW_DASHBOARD = "/dashboard"
	PATH_ROOT           = "/"
	PATH_CSS            = "/styles.css"
	PATH_SYSTEMJS_CFG   = "/systemjs.config.js"

	TF_CLF string = "02/Jan/2006:15:04:05 -0700"
)

type WebResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WebWhitelistAddRequest struct {
	IpAddr      string `json:"ip_addr"`
	Description string `json:"description"`
}

type WebWhitelistRemoveRequest struct {
	IpAddr string `json:"ip_addr"`
}

type WebESProxyDetailsRequest struct {
	FlowId int64 `json:"flow_id"`
}

type TemplateData struct{}

type RtbhApi struct {
	cfg       *config.Config
	log       *logger.Logger
	blacklist *blacklist.Blacklist
	whitelist *whitelist.Whitelist
	history   *history.History
	es        *es.ES
}
