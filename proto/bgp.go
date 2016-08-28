package proto

import (
	"errors"
	"github.com/r3boot/go-rtbh/proto/bgp2go"
	"net"
	"strconv"
	"strings"
)

var BGPContext bgp2go.BGPContext

var cmdToPeer chan bgp2go.BGPProcessMsg
var cmdFromPeer chan bgp2go.BGPProcessMsg

// Utility function used to convert a router id string to uint32
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

// Utility function used to convert a community string to uint32
func community_aton(community string) uint32 {
	var sum uint32
	var words []string
	var w0 int
	var w1 int

	words = strings.Split(community, ":")
	w0, _ = strconv.Atoi(words[0])
	w1, _ = strconv.Atoi(words[1])

	sum += uint32(w0) << 16
	sum += uint32(w1)

	return sum
}

// Utility function used to add a default prefixlen to a prefix if needed
func add_cidr_mask(addr string) string {
	if strings.Contains(addr, "/") {
		return addr
	}

	if strings.Contains(addr, ":") {
		return addr + "/128"
	} else {
		return addr + "/32"
	}
}

// Configure this side of the BGP routine
func ConfigureBGP() (err error) {
	var asnum uint64

	// Convert asnum string into uint64 for later use
	if asnum, err = strconv.ParseUint(Config.BGP.Asnum, 10, 32); err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse Asnum: " + err.Error())
		return
	}

	BGPContext.ASN = uint32(asnum)
	BGPContext.ListenLocal = true

	BGPContext.RouterID, err = bgp2go.IPv4ToUint32(Config.BGP.RouterId)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse RouterID: " + err.Error())
		return
	}

	BGPContext.NextHop, err = bgp2go.IPv4ToUint32(Config.BGP.NextHop)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse IPv4 NextHop: " + err.Error())
		return
	}

	BGPContext.NextHopV6, err = bgp2go.IPv6StringToAddr(Config.BGP.NextHopV6)
	if err != nil {
		err = errors.New("[ConfigureBGP]: Failed to parse IPv6 NextHop: " + err.Error())
		return
	}

	BGPContext.Community = append(BGPContext.Community, community_aton(Config.BGP.Community))
	BGPContext.LocalPref = uint32(Config.BGP.LocalPref)

	return
}

func RunBGP() {
	cmdToPeer = make(chan bgp2go.BGPProcessMsg)
	cmdFromPeer = make(chan bgp2go.BGPProcessMsg)

	Log.Debug("Starting BGP routine")
	go bgp2go.StartBGPProcess(cmdToPeer, cmdFromPeer, BGPContext)
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
	prefix = add_cidr_mask(prefix)
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
	prefix = add_cidr_mask(prefix)
	if strings.Contains(prefix, ":") {
		RemoveBGPv6Route(prefix)
	} else {
		RemoveBGPv4Route(prefix)
	}
}
