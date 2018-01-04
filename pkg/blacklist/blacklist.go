package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func New(l *logger.Logger, c *config.Config, b *bgp.BGP) *Blacklist {
	log = l
	cfg = c

	bl := &Blacklist{
		bgp:   b,
		mutex: &sync.Mutex{},
	}

	log.Debugf("Blacklist: Module initialized")

	return bl
}
