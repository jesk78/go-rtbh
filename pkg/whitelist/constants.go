package whitelist

import (
	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

type Whitelist struct {
	bgp *bgp.BGP
}

var (
	cfg *config.Config
	log *logger.Logger
)
