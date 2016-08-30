package lib

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/go-rtbh/lib/reaper"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/lib/whitelist"
	"github.com/r3boot/rlib/logger"
)

var Config config.Config
var Log logger.Log

var Blacklist *blacklist.Blacklist
var Whitelist *whitelist.Whitelist
var History *history.History
var Resolver *resolver.Resolver
var Reaper *reaper.Reaper

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	if err = listcache.Setup(Log, Config); err != nil {
		return
	}

	if err = blacklist.Setup(Log, Config); err != nil {
		return
	}
	Blacklist = blacklist.New()

	if err = whitelist.Setup(Log, Config); err != nil {
		return
	}
	Whitelist = whitelist.New()

	if err = history.Setup(Log, Config); err != nil {
		return
	}
	History = history.New()

	if err = resolver.Setup(Log, Config); err != nil {
		return
	}
	Resolver = resolver.New()

	if err = reaper.Setup(Log, Config); err != nil {
		return
	}
	Reaper = reaper.New(Blacklist)

	Log.Debug("Lib: All submodules initialized")
	return
}
