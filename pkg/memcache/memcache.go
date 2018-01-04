package memcache

import (
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/rlib/logger"
	"sync"
)

const MYNAME string = "ListCache"

var Config *config.Config
var Log logger.Log

type StringCache struct {
	cache map[string]interface{}
	mutex *sync.Mutex
}

type IntCache struct {
	cache map[int64]interface{}
	mutex *sync.Mutex
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func NewStringIndexed() *StringCache {
	var cache *StringCache

	cache = &StringCache{
		cache: make(map[string]interface{}),
		mutex: &sync.Mutex{},
	}

	return cache
}

func NewIntIndexed() *IntCache {
	var cache *IntCache

	cache = &IntCache{
		cache: make(map[int64]interface{}),
		mutex: &sync.Mutex{},
	}

	return cache
}
