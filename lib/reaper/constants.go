package reaper

import (
	"time"

	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
)

type Reaper struct {
	Interval  time.Duration
	Control   chan int
	Done      chan bool
	blacklist *blacklist.Blacklist
}

var (
	cfg *config.Config
	log *logger.Logger
)
