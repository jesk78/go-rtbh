package lib

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/amqp"
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/pipeline"
	"github.com/r3boot/go-rtbh/lib/reaper"
	"github.com/r3boot/go-rtbh/lib/redis"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/lib/whitelist"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "Lib"

var Config config.Config
var Log logger.Log

var AmqpClient *amqp.AmqpClient
var RedisClient *redis.RedisClient
var Blacklist *blacklist.Blacklist
var Whitelist *whitelist.Whitelist
var History *history.History
var Resolver *resolver.Resolver
var Reaper *reaper.Reaper
var Pipeline *pipeline.Pipeline
var BGP *bgp.BGP
var ORM *orm.ORM

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	// First, configure all dependencies
	if err = config.Setup(Log); err != nil {
		return
	}
	Config = config.New(*cfgfile)

	if err = events.Setup(Log, Config); err != nil {
		return
	}

	if Config.Redis.Address == "" && Config.Amqp.Address == "" {
		err = errors.New(MYNAME + ": No event feed to connect to")
	}

	if Config.Amqp.Address != "" {
		if err = amqp.Setup(Log, Config); err != nil {
			return
		}
		AmqpClient = amqp.New()
	}

	if Config.Redis.Address != "" {
		if err = redis.Setup(Log, Config); err != nil {
			return
		}
		RedisClient = redis.New()
	}

	if err = orm.Setup(Log, Config); err != nil {
		return
	}
	ORM = orm.New()

	if err = bgp.Setup(Log, Config); err != nil {
		return
	}
	BGP = bgp.New()

	// Then, setup all blacklist related libs
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

	if err = pipeline.Setup(Log, Config); err != nil {
		return
	}
	Pipeline = pipeline.New(Config.Ruleset, Blacklist, Whitelist, History)

	Log.Debug("Lib: All submodules initialized")
	return
}
