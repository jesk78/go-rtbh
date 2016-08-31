package blacklist

import (
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/rlib/logger"
	"sync"
)

const MYNAME string = "Blacklist"

type Blacklist struct {
	cache *listcache.Cache
	bgp   *bgp.BGP
	mutex *sync.Mutex
}

var Config *config.Config
var Log logger.Log

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": initialized with entries")
	return
}

func New(b *bgp.BGP) *Blacklist {
	var bl *Blacklist

	bl = &Blacklist{
		cache: listcache.New(),
		bgp:   b,
		mutex: &sync.Mutex{},
	}

	return bl
}
