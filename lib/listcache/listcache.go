package listcache

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/logger"
	"sync"
)

const MYNAME string = "ListCache"

var Config config.Config
var Log logger.Log

type Cache struct {
	cache map[string]interface{}
	mutex *sync.Mutex
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New() *Cache {
	var cache *Cache

	cache = &Cache{
		cache: make(map[string]interface{}),
		mutex: &sync.Mutex{},
	}

	return cache
}
