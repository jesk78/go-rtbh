package whitelist

import (
	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func NewWhitelist(l *logger.Logger, c *config.Config, b *bgp.BGP) *Whitelist {
	log = l
	cfg = c

	wl := &Whitelist{
		bgp: b,
	}

	log.Debugf("Whitelist: Module initialized")

	return wl
}
