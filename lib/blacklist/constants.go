package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/logger"
)

type ApiBlacklistGetAllResponse struct {
	data []*events.APIEvent
}

type Blacklist struct {
	bgp   *bgp.BGP
	mutex *sync.Mutex
}

var (
	cfg *config.Config
	log *logger.Logger
)
