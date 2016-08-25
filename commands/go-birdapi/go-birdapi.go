package main

import (
	"flag"
	"github.com/r3boot/go-rtbh/birdapi"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
	"os"
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

	if err = birdapi.Setup(Log, Config, Redis); err != nil {
		Log.Fatal(err)
	}
}

func main() {
	birdapi.RunServer()

	// Wait for program completion
	Log.Debug("[go-birdapi]: Program finished")
}
