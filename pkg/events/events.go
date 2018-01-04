package events

import (
	"time"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

const MYNAME string = "Events"

var (
	cfg *config.Config
	log *logger.Logger
)

func Setup(l *logger.Logger, c *config.Config) {
	log = l
	cfg = c

	log.Debugf("Events: Module initialized")
	return
}

func New(data []byte) (event *RTBHEvent, err error) {
	event = &RTBHEvent{
		AddedAt: time.Now(),
	}
	err = event.LoadFrom(data)
	if err != nil {
		event = nil
	}

	return
}
