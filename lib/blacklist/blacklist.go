package blacklist

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "Blacklist"

type Blacklist struct {
	cache *listcache.Cache
}

var Config config.Config
var Log logger.Log

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": initialized with entries")
	return
}

func New() *Blacklist {
	var bl *Blacklist

	bl = &Blacklist{
		cache: listcache.New(),
	}

	return bl
}
