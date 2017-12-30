package events

import (
	"time"

	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
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

func NewEvent(data []byte) (event *RTBHEvent, err error) {
	event = &RTBHEvent{
		AddedAt: time.Now(),
	}
	err = event.LoadFrom(data)
	if err != nil {
		event = nil
	}

	return
}
