package pipeline

import (
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/logger"
	"github.com/r3boot/go-rtbh/lib/whitelist"
)

type Pipeline struct {
	blacklist *blacklist.Blacklist
	whitelist *whitelist.Whitelist
	history   *history.History
	Control   chan int
	Done      chan bool
}

var (
	cfg *config.Config
	log *logger.Logger
)
