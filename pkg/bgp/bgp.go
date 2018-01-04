package bgp

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/bgp/bgp2go"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func NewBGP(l *logger.Logger, c *config.Config) (*BGP, error) {
	log = l
	cfg = c

	bgp := &BGP{
		context:     bgp2go.BGPContext{},
		cmdToPeer:   make(chan bgp2go.BGPProcessMsg),
		cmdFromPeer: make(chan bgp2go.BGPProcessMsg),
	}

	err := bgp.Configure()
	if err != nil {
		return nil, fmt.Errorf("NewBGP: %v", err)
	}

	log.Debugf("BGP: Module initialized")

	return bgp, nil
}
