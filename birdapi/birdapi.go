package birdapi

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
)

var Bird *proto.Bird
var Config *config.Config
var Log logger.Log
var Redis *redis.Client

var Blacklist *lists.Blacklist
var Whitelist *lists.Whitelist

func Setup(l logger.Log, cfg *config.Config, r *redis.Client) (err error) {
	Log = l
	Config = cfg
	Redis = r
	Bird = proto.NewBirdClient()

	Blacklist = lists.NewBlacklist()
	Whitelist = lists.NewWhitelist()

	ConfigureRouting()

	return
}
