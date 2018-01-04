package reaper

import (
	"time"

	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
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
