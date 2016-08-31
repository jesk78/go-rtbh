package bgp

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lib/bgp/bgp2go"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "BGP"

var Config config.Config
var Log logger.Log

type BGP struct {
	context     bgp2go.BGPContext
	cmdToPeer   chan bgp2go.BGPProcessMsg
	cmdFromPeer chan bgp2go.BGPProcessMsg
}

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}

func New() *BGP {
	var bgp *BGP
	var err error

	bgp = &BGP{}
	if err = bgp.Configure(); err != nil {
		return nil
	}

	return bgp
}
