package resolver

import (
	"fmt"
	"sync"
	"time"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func New(l *logger.Logger, c *config.Config, o *orm.ORM) (resolver *Resolver, err error) {
	if !c.General.Resolver.Enabled {
		return nil, fmt.Errorf("NewResolver: DNS resolver not enabled")
	}

	resolver = &Resolver{
		log:     l,
		cfg:     c,
		orm:     o,
		cache:   make(map[string]string),
		mutex:   &sync.Mutex{},
		Control: make(chan int, config.D_CONTROL_BUFSIZE),
		Done:    make(chan bool, config.D_DONE_BUFSIZE),
	}

	resolver.Interval, err = time.ParseDuration(resolver.cfg.General.Resolver.LookupMaxInterval)
	if err != nil {
		return nil, fmt.Errorf("NewResolver time.ParseDuration: %v", err)
	}

	resolver.log.Debugf("Resolver: Module initialized")

	return
}
