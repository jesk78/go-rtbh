package pipeline

import (
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
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
