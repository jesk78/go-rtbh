package bgp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bgp2go"

	"github.com/r3boot/go-rtbh/pkg/config"
)

// Configure this side of the BGP routine
func (bgp *BGP) Configure() error {
	// Convert asnum string into uint64 for later use
	asnum, err := strconv.ParseUint(cfg.BGP.Asnum, 10, 32)
	if err != nil {
		return fmt.Errorf("BGP.Configure strconv.ParseUint: %v", err)

	}

	bgp.context.ASN = uint32(asnum)
	bgp.context.ListenLocal = true

	bgp.context.RouterID, err = bgp2go.IPv4ToUint32(cfg.BGP.RouterId)
	if err != nil {
		return fmt.Errorf("BGP.Configure bgp2go.IPv4ToUint32: %v", err)
	}

	bgp.context.NextHop, err = bgp2go.IPv4ToUint32(cfg.BGP.NextHop)
	if err != nil {
		return fmt.Errorf("BGP.Configure bgp2go.IPv4ToUint32: %v", err)
	}

	bgp.context.NextHopV6, err = bgp2go.IPv6StringToAddr(cfg.BGP.NextHopV6)
	if err != nil {
		return fmt.Errorf("BGP.Configure bgp2go.IPv6StringToAddr: %v", err)
	}

	bgp.context.Community = append(bgp.context.Community, community_aton(cfg.BGP.Community))
	bgp.context.LocalPref = uint32(cfg.BGP.LocalPref)

	return nil
}

func (bgp *BGP) ServerRoutine() {
	var bgpPeer config.BGPPeer

	bgp.cmdToPeer = make(chan bgp2go.BGPProcessMsg)
	bgp.cmdFromPeer = make(chan bgp2go.BGPProcessMsg)

	log.Debugf("BGP.ServerRoutine: Starting bgp process")
	go bgp2go.StartBGPProcess(bgp.cmdToPeer, bgp.cmdFromPeer, bgp.context)

	time.Sleep(1 * time.Second)
	for _, bgpPeer = range cfg.BGP.Peers {
		bgp.AddNeighbor(bgpPeer.Address)
	}
}

func (bgp *BGP) addv4Neighbor(ipaddr string) {
	log.Debugf("BGP.addv4Neighbor: Adding %s", ipaddr)
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddNeighbour",
		Data: ipaddr + " inet",
	}
}

func (bgp *BGP) addv6Neighbor(ipaddr string) {
	log.Debugf("BGP.addv6Neighbor: Adding %s", ipaddr)
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
	log.Debugf("BGP.removev4Neighbor: Removing %s", ipaddr)
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "RemoveNeighbour",
		Data: ipaddr + " inet",
	}
}

func (bgp *BGP) removev6Neighbor(ipaddr string) {
	log.Debugf("BGP.removev6Neighbor: Removing %s", ipaddr)
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
	log.Debugf("BGP.addv4Route: Adding %s", prefix)
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddV4Route",
		Data: prefix,
	}
}

func (bgp *BGP) addv6Route(prefix string) {
	log.Debugf("BGP.addv6Route: Adding %s", prefix)
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
	log.Debugf("BGP.removev4Route: Removing %s", prefix)
	bgp.cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "WithdrawV4Route",
		Data: prefix,
	}
}

func (bgp *BGP) removev6Route(prefix string) {
	log.Debugf("BGP.removev6Route: Removing %s", prefix)
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
