package lists

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
)

var Log logger.Log
var Config *config.Config
var Redis *redis.Client

func Setup(l logger.Log, cfg *config.Config, rc *redis.Client) (err error) {
	Log = l
	Config = cfg
	Redis = rc

	return
}
