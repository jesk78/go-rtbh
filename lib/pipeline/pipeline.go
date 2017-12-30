package pipeline

import (
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/logger"
	"github.com/r3boot/go-rtbh/lib/whitelist"
)

func NewPipeline(l *logger.Logger, c *config.Config, bl *blacklist.Blacklist, wl *whitelist.Whitelist, h *history.History) *Pipeline {
	log = l
	cfg = c

	pipeline := &Pipeline{
		blacklist: bl,
		whitelist: wl,
		history:   h,
		Control:   make(chan int, config.D_CONTROL_BUFSIZE),
		Done:      make(chan bool, config.D_DONE_BUFSIZE),
	}

	log.Debugf("Pipeline: Module initialized")

	return pipeline
}
