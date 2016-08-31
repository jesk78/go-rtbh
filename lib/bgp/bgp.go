package bgp

import (
	"github.com/r3boot/go-rtbh/lib/bgp/bgp2go"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "BGP"

var Config *config.Config
var Log logger.Log

type BGP struct {
	context     bgp2go.BGPContext
	cmdToPeer   chan bgp2go.BGPProcessMsg
	cmdFromPeer chan bgp2go.BGPProcessMsg
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New() BGP {
	var bgp BGP
	var err error

	bgp = BGP{
		context:     bgp2go.BGPContext{},
		cmdToPeer:   make(chan bgp2go.BGPProcessMsg),
		cmdFromPeer: make(chan bgp2go.BGPProcessMsg),
	}
	if err = bgp.Configure(); err != nil {
		Log.Fatal(err)
	}

	return bgp
}
