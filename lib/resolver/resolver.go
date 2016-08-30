package resolver

import (
	"errors"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/logger"
	"sync"
	"time"
)

const MYNAME string = "Resolver"

var Config config.Config
var Log logger.Log

type Resolver struct {
	Interval time.Duration
	Control  chan int
	Done     chan bool
	cache    map[string]string
	mutex    *sync.Mutex
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug("Resolver: DNS resolver initialized")

	return
}

func NewResolver() (resolver *Resolver, err error) {
	if !Config.General.Resolver.Enabled {
		err = errors.New(MYNAME + ": dns resolving not enabled")
		return
	}

	resolver = &Resolver{}
	resolver.Interval, err = time.ParseDuration(Config.Resolver.LookupMaxInterval)
	if err != nil {
		resolver = nil
		err = errors.New(MYNAME + ": Failed to parse duration: " + err.Error())
		return
	}

	resolver.cache = make(map[string]string)
	resolver.mutex = &sync.Mutex{}

	Log.Debug(MYNAME + ": Initialized new dns resolver")

	return
}
