package bgp

import (
	"github.com/r3boot/go-rtbh/lib/bgp/bgp2go"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/logger"
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
