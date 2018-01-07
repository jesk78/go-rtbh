package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

type ApiBlacklistGetAllResponse struct {
	data []*events.APIEvent
}

type Blacklist struct {
	cfg   *config.Config
	log   *logger.Logger
	orm   *orm.ORM
	bgp   *bgp.BGP
	mutex *sync.Mutex
}
