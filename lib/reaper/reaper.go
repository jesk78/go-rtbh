package reaper

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
	"time"
)

const MYNAME string = "Reaper"

var Config *config.Config
var Log logger.Log

type Reaper struct {
	Interval  time.Duration
	Control   chan int
	Done      chan bool
	blacklist *blacklist.Blacklist
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New(bl *blacklist.Blacklist) *Reaper {
	var reaper *Reaper
	var err error

	reaper = &Reaper{
		blacklist: bl,
	}

	reaper.Interval, err = time.ParseDuration(Config.General.ReaperInterval)
	if err != nil {
		err = errors.New(MYNAME + ": Failed to parse duration: " + err.Error())
		return nil
	}

	return reaper
}
