package proto

import (
	"errors"
	"github.com/tehnerd/bgp2go"
	"net"
	"strconv"
	"strings"
)

var BGPContext bgp2go.BGPContext

var cmdToPeer chan bgp2go.BGPProcessMsg
var cmdFromPeer chan bgp2go.BGPProcessMsg

func inet_aton(ip net.IP) uint32 {
	var bits []string
	var b0 int
	var b1 int
	var b2 int
	var b3 int
	var sum uint32

	bits = strings.Split(ip.String(), ".")
	b0, _ = strconv.Atoi(bits[0])
	b1, _ = strconv.Atoi(bits[1])
	b2, _ = strconv.Atoi(bits[2])
	b3, _ = strconv.Atoi(bits[3])

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

func ConfigureBGPd() (err error) {
	var asnum uint64
	var routerid_ip net.IP

	// Convert asnum string into uint64 for later use
	if asnum, err = strconv.ParseUint(Config.BGP.Asnum, 10, 32); err != nil {
		err = errors.New("[ConfigureBGPd]: Failed to parse ASNum: " + err.Error())
		return
	}

	// Convert routerid string into uint32
	routerid_ip = net.ParseIP(Config.BGP.RouterId)
	if routerid_ip == nil {
		err = errors.New("[ConfigureBGPd]: Failed to convert RouterID to uint32")
		return
	}

	BGPContext.ASN = uint32(asnum)
	BGPContext.RouterID = inet_aton(routerid_ip)
	BGPContext.ListenLocal = true

	return
}

func RunBGPd() {
	cmdToPeer = make(chan bgp2go.BGPProcessMsg)
	cmdFromPeer = make(chan bgp2go.BGPProcessMsg)

	Log.Debug("Starting BGP routine")
	bgp2go.StartBGPProcess(cmdToPeer, cmdFromPeer, BGPContext)
}

func AddBGPv4Neighbor(ipaddr string) {
	Log.Debug("Adding IPv4 BGP neighbor")
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddNeighbour",
		Data: ipaddr + " inet",
	}
}

func AddBGPv6Neighbor(ipaddr string) {
	Log.Debug("Adding IPv6 BGP neighbor")
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddNeighbour",
		Data: ipaddr + " inet6",
	}
}

func AddBGPNeighbor(ipaddr string) {
	if strings.Contains(ipaddr, ":") {
		AddBGPv6Neighbor(ipaddr)
	} else {
		AddBGPv4Neighbor(ipaddr)
	}
}

func RemoveBGPv4Neighbor(ipaddr string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "RemoveNeighbour",
		Data: ipaddr + " inet",
	}
}

func RemoveBGPv6Neighbor(ipaddr string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "RemoveNeighbour",
		Data: ipaddr + " inet6",
	}
}

func RemoveBGPNeighbor(ipaddr string) {
	if strings.Contains(ipaddr, ":") {
		RemoveBGPv6Neighbor(ipaddr)
	} else {
		RemoveBGPv4Neighbor(ipaddr)
	}
}

func AddBGPv4Route(prefix string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddV4Route",
		Data: prefix,
	}
}

func AddBGPv6Route(prefix string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "AddV6Route",
		Data: prefix,
	}
}

func AddBGPRoute(prefix string) {
	if strings.Contains(prefix, ":") {
		AddBGPv6Route(prefix)
	} else {
		AddBGPv4Route(prefix)
	}
}

func RemoveBGPv4Route(prefix string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "WithdrawV4Route",
		Data: prefix,
	}
}

func RemoveBGPv6Route(prefix string) {
	cmdToPeer <- bgp2go.BGPProcessMsg{
		Cmnd: "WithdrawV6Route",
		Data: prefix,
	}
}

func RemoveBGPRoute(prefix string) {
	if strings.Contains(prefix, ":") {
		RemoveBGPv6Route(prefix)
	} else {
		RemoveBGPv4Route(prefix)
	}
}
