package bgp

import (
	"errors"
	"github.com/r3boot/go-rtbh/lib/bgp/bgp2go"
	"strconv"
	"strings"
)

// Configure this side of the BGP routine
func (bgp *BGP) Configure() (err error) {
	var asnum uint64

	// Convert asnum string into uint64 for later use
	if asnum, err = strconv.ParseUint(Config.BGP.Asnum, 10, 32); err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse Asnum: " + err.Error())
		return
	}

	bgp.context.ASN = uint32(asnum)
	bgp.context.ListenLocal = true

	bgp.context.RouterID, err = bgp2go.IPv4ToUint32(Config.BGP.RouterId)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse RouterID: " + err.Error())
		return
	}

	bgp.context.NextHop, err = bgp2go.IPv4ToUint32(Config.BGP.NextHop)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse IPv4 NextHop: " + err.Error())
		return
	}

	bgp.context.NextHopV6, err = bgp2go.IPv6StringToAddr(Config.BGP.NextHopV6)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse IPv6 NextHop: " + err.Error())
		return
	}

	bgp.context.Community = append(bgp.context.Community, community_aton(Config.BGP.Community))
	bgp.context.LocalPref = uint32(Config.BGP.LocalPref)

	return
}

func (bgp *BGP) ServerRoutine() {
	bgp.cmdToPeer = make(chan bgp2go.BGPProcessMsg)
	bgp.cmdFromPeer = make(chan bgp2go.BGPProcessMsg)

	Log.Debug("Starting BGP routine")
	go bgp2go.StartBGPProcess(bgp.cmdToPeer, bgp.cmdFromPeer, bgp.context)
}

func (bgp *BGP) addv4Neighbor(ipaddr string) {
	Log.Debug("Adding IPv4 BGP neighbor")
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddNeighbour",
		Data: ipaddr + " inet",
	}
}

func (bgp *BGP) addv6Neighbor(ipaddr string) {
	Log.Debug("Adding IPv6 BGP neighbor")
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddNeighbour",
		Data: ipaddr + " inet6",
	}
}

func (bgp *BGP) AddNeighbor(ipaddr string) {
	if strings.Contains(ipaddr, ":") {
		bgp.addv6Neighbor(ipaddr)
	} else {
		bgp.addv4Neighbor(ipaddr)
	}
}

func (bgp *BGP) removev4Neighbor(ipaddr string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "RemoveNeighbour",
		Data: ipaddr + " inet",
	}
}

func (bgp *BGP) removev6Neighbor(ipaddr string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "RemoveNeighbour",
		Data: ipaddr + " inet6",
	}
}

func (bgp *BGP) RemoveBGPNeighbor(ipaddr string) {
	if strings.Contains(ipaddr, ":") {
		bgp.removev6Neighbor(ipaddr)
	} else {
		bgp.removev4Neighbor(ipaddr)
	}
}

func (bgp *BGP) addv4Route(prefix string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddV4Route",
		Data: prefix,
	}
}

func (bgp *BGP) addv6Route(prefix string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddV6Route",
		Data: prefix,
	}
}

func (bgp *BGP) AddRoute(prefix string) {
	prefix = add_cidr_mask(prefix)
	if strings.Contains(prefix, ":") {
		bgp.addv6Route(prefix)
	} else {
		bgp.addv4Route(prefix)
	}
}

func (bgp *BGP) removev4Route(prefix string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "WithdrawV4Route",
		Data: prefix,
	}
}

func (bgp *BGP) removev6Route(prefix string) {
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "WithdrawV6Route",
		Data: prefix,
	}
}

func (bgp *BGP) RemoveRoute(prefix string) {
	prefix = add_cidr_mask(prefix)
	if strings.Contains(prefix, ":") {
		bgp.removev6Route(prefix)
	} else {
		bgp.removev4Route(prefix)
	}
}
