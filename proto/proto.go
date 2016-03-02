package proto

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/logger"
)

var Log logger.Log
var Config *config.Config

func Setup(l logger.Log, cfg *config.Config) (err error) {
	Log = l
	Config = cfg

	return
}
