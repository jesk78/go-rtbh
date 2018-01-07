package whitelist

import (
	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func New(l *logger.Logger, c *config.Config, o *orm.ORM, b *bgp.BGP) *Whitelist {

	wl := &Whitelist{
		log: l,
		cfg: c,
		orm: o,
		bgp: b,
	}

	wl.log.Debugf("Whitelist: Module initialized")

	return wl
}
