package events

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
	"time"
)

const MYNAME string = "Events"

var Config config.Config
var Log logger.Log

type RTBHEvent struct {
	Address  string
	Reason   string
	AddedAt  time.Time
	ExpireIn string
}

// Struct containing a whitelist entry
type RTBHWhiteEntry struct {
	Address     string
	Description string
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

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
