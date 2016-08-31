package history

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "History"

var Config *config.Config
var Log logger.Log

type History struct {
	cache *listcache.Cache
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New() *History {
	var history *History

	history = &History{
		cache: listcache.New(),
	}

	return history
}
