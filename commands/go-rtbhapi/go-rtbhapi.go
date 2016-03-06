package main

import (
	"flag"
	"github.com/r3boot/go-rtbh/api"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
)

var Config *config.Config
var Log logger.Log
var Redis *redis.Client

var cfgfile = flag.String("f", config.D_APICFGFILE, "Configuration file to use")
var debug = flag.Bool("D", config.D_DEBUG, "Enable debug output")
var timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamps

	// Setup config library
	if err = config.Setup(Log); err != nil {
		Log.Fatal("[config]: Initialization failed: " + err.Error())
	}
	Log.Debug("[config]: Initialization succesful")

	// Fetch a new configuration struct
	if Config, err = config.NewConfig(); err != nil {
		Log.Fatal("[config]: Failed to generate empty config")
	}

	if err = Config.LoadFrom(*cfgfile); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("[config]: Loaded configuration from " + *cfgfile)

	// Prepare protocols for usage
	if err = proto.Setup(Log, Config); err != nil {
		Log.Fatal("[proto]: Initialization failed: " + err.Error())
	}
	Log.Debug("[proto]: Initialization succesful")

	// Setup redis client
	if Redis, err = proto.NewRedisClient(); err != nil {
		Log.Fatal("[redis]: Initialization failed: " + err.Error())
	}
	Log.Debug("[redis]: Connected to " + Config.Redis.Address)

	// Prepare lists for usage
	if err = lists.Setup(Log, Config, Redis); err != nil {
		Log.Fatal("[lists]: Initialization failed: " + err.Error())
	}
	Log.Debug("[lists]: Initialization succesful")

	if err = api.Setup(Log, Config); err != nil {
		Log.Fatal("[api]: Initialization failed: " + err.Error())
	}
	Log.Debug("[api]: Initialization succesful")
}

func main() {
	// Fire up the api
	api.RunTillDeath()
}
