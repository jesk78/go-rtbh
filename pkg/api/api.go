package api

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/es"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
)

func New(l *logger.Logger, c *config.Config, b *blacklist.Blacklist, w *whitelist.Whitelist, h *history.History, e *es.ES) *RtbhApi {
	api := &RtbhApi{
		log:       l,
		cfg:       c,
		blacklist: b,
		whitelist: w,
		history:   h,
		es:        e,
	}

	return api
}
