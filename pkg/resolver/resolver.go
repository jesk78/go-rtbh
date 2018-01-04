package resolver

import (
	"fmt"
	"sync"
	"time"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func NewResolver(l *logger.Logger, c *config.Config) (resolver *Resolver, err error) {
	log = l
	cfg = c

	if !cfg.General.Resolver.Enabled {
		return nil, fmt.Errorf("NewResolver: DNS resolver not enabled")
	}

	resolver = &Resolver{
		cache:   make(map[string]string),
		mutex:   &sync.Mutex{},
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	resolver.Interval, err = time.ParseDuration(cfg.General.Resolver.LookupMaxInterval)
	if err != nil {
		return nil, fmt.Errorf("NewResolver time.ParseDuration: %v", err)
	}

	log.Debugf("Resolver: Module initialized")

	return
}
