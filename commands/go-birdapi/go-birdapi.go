package main

import (
	"flag"
	_ "github.com/r3boot/go-rtbh/birdapi"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
	"os"
	"time"
)

// Program libraries
var Log logger.Log
var Config *config.Config
var Bird *proto.Bird
var Redis *redis.Client

// Lists
var Blacklist *lists.Blacklist
var Whitelist *lists.Whitelist

// OS signals
var signals chan os.Signal
var allDone chan bool

// Command-line flags
var cfgfile = flag.String("f", config.D_CFGFILE, "Configuration file to use")
var debug = flag.Bool("D", config.D_DEBUG, "Enable debug output")
var timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamps
	Log.Debug("Logging initialized")

	// Setup configuration library
	if err = config.Setup(Log); err != nil {
		Log.Fatal("[config]: Initialization failed: " + err.Error())
	}
	Log.Debug("[config]: Library initialized")

	// Generate configuration struct and try to load configuration
	if Config, err = config.NewConfig(); err != nil {
		Log.Fatal("[config]: Failed to create empty configuration: " + err.Error())
	}

	if err = Config.LoadFrom(*cfgfile); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("Loaded configuration from " + *cfgfile)

	// Setup protocol library
	if err = proto.Setup(Log, Config); err != nil {
		Log.Fatal("[proto]: Initialization failed: " + err.Error())
	}
	Log.Debug("[proto]: Library initialized")

	// Setup redis client
	if Redis, err = proto.NewRedisClient(); err != nil {
		Log.Fatal("[Redis]: Failed to initialize: " + err.Error())
	}
	Log.Debug("[Redis]: Connected to " + Config.Redis.Address)

	// Setup lists library
	if err = lists.Setup(Log, Config, Redis); err != nil {
		Log.Fatal("[lists]: Initialization failed: " + err.Error())
	}
	Log.Debug("[lists]: Library initialized")

	// Initialize BGP daemon
	if err = proto.ConfigureBGPd(); err != nil {
		Log.Fatal("[BGP]: Initialization failed: " + err.Error())
	}
	Log.Debug("[BGP]: Library initialized")

}

func main() {
	var bgpPeer config.BGPPeer

	go proto.RunBGPd()

	time.Sleep(1 * time.Second)

	for _, bgpPeer = range Config.BGP.Peers {
		Log.Debug("Adding BGP neighbor " + bgpPeer.Address)
		proto.AddBGPNeighbor(bgpPeer.Address)
	}

	for _, entry := range Blacklist.GetAll() {
		Log.Debug("Adding BGP route " + entry)
		proto.AddBGPRoute(entry + "/32")
	}

	Log.Debug(proto.BGPContext.Neighbours)
	Log.Debug(proto.BGPContext.RIBv4)

	time.Sleep(60 * time.Second)

	// Wait for program completion
	Log.Debug("[go-birdapi]: Program finished")
}
