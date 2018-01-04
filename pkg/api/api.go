package api

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "API"

var Config *config.Config
var Log logger.Log

type RtbhApi struct {
	blacklist *blacklist.Blacklist
	whitelist *whitelist.Whitelist
	history   *history.History
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New(b *blacklist.Blacklist, w *whitelist.Whitelist, h *history.History) *RtbhApi {
	var api *RtbhApi

	api = &RtbhApi{
		blacklist: b,
		whitelist: w,
		history:   h,
	}

	return api
}
