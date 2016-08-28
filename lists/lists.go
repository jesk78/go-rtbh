package lists

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/logger"
)

const REDIS_BASE string = "net.as65342.go-rtbh"
const REDIS_ADDR_SEPARATOR string = "^"

var Log logger.Log
var Config *config.Config

func Setup(l logger.Log, cfg *config.Config) (err error) {
	Log = l
	Config = cfg

	return
}
