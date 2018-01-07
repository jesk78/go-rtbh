package whitelist

import (
	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

type Whitelist struct {
	cfg *config.Config
	log *logger.Logger
	orm *orm.ORM
	bgp *bgp.BGP
}
