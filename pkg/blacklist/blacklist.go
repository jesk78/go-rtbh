package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func New(l *logger.Logger, c *config.Config, o *orm.ORM, b *bgp.BGP) *Blacklist {
	bl := &Blacklist{
		cfg:   c,
		log:   l,
		orm:   o,
		bgp:   b,
		mutex: &sync.Mutex{},
	}

	bl.log.Debugf("Blacklist: Module initialized")

	return bl
}
