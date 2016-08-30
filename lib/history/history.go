package history

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "History"

var Config config.Config
var Log logger.Log

type History struct {
	cache *listcache.Cache
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New() *History {
	var history *History

	history = &History{
		cache: listcache.New(),
	}

	return history
}
