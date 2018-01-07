package api

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
)

func New(l *logger.Logger, c *config.Config, b *blacklist.Blacklist, w *whitelist.Whitelist, h *history.History) *RtbhApi {
	log = l
	cfg = c

	api := &RtbhApi{
		blacklist: b,
		whitelist: w,
		history:   h,
	}

	return api
}
