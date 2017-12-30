package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
)

func NewBlacklist(l *logger.Logger, c *config.Config, b *bgp.BGP) *Blacklist {
	log = l
	cfg = c

	bl := &Blacklist{
		bgp:   b,
		mutex: &sync.Mutex{},
	}

	log.Debugf("Blacklist: Module initialized")

	return bl
}
