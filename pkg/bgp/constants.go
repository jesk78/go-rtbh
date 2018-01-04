package bgp

import (
	"bgp2go"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

type BGP struct {
	context     bgp2go.BGPContext
	cmdToPeer   chan bgp2go.BGPProcessMsg
	cmdFromPeer chan bgp2go.BGPProcessMsg
}

var (
	cfg *config.Config
	log *logger.Logger
)
