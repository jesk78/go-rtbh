package reaper

import (
	"fmt"
	"time"

	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func New(l *logger.Logger, c *config.Config, bl *blacklist.Blacklist) (*Reaper, error) {
	var err error

	log = l
	cfg = c

	reaper := &Reaper{
		blacklist: bl,
	}

	reaper.Interval, err = time.ParseDuration(cfg.General.ReaperInterval)
	if err != nil {
		return nil, fmt.Errorf("NewReaper time.ParseDuration: %v", err)
	}

	log.Debugf("Reaper: Module initialized")

	return reaper, nil
}
