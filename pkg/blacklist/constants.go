package blacklist

import (
	"sync"

	"github.com/r3boot/go-rtbh/pkg/bgp"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/logger"
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
