package pipeline

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
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
