package whitelist

import (
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "Whitelist"

type Whitelist struct {
	cache *listcache.Cache
	bgp   *bgp.BGP
}

var Config *config.Config
var Log logger.Log

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New(b *bgp.BGP) *Whitelist {
	var wl *Whitelist

	wl = &Whitelist{
		cache: listcache.New(),
		bgp:   b,
	}

	return wl
}
