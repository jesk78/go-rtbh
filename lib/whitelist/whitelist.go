package whitelist

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "WHITELIST"

type Whitelist struct {
	cache *listcache.Cache
}

var Config config.Config
var Log logger.Log

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New() *Whitelist {
	var wl *Whitelist

	wl = &Whitelist{
		cache: listcache.New(),
	}

	return wl
}
