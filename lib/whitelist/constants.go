package whitelist

import (
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
)

type Whitelist struct {
	bgp *bgp.BGP
}

var (
	cfg *config.Config
	log *logger.Logger
)
