package main

import (
	"flag"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/pipeline"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
	"os"
	"os/signal"
	"strconv"
)

// Program libraries
var Log logger.Log
var Config *config.Config
var Amqp *proto.AmqpClient
var Redis *redis.Client
var Pipeline *pipeline.Pipeline
var Whitelist *lists.Whitelist

// OS signals
var signals chan os.Signal
var allDone chan bool

// Command-line flags
var cfgfile = flag.String("f", config.D_CFGFILE, "Configuration file to use")
var debug = flag.Bool("D", config.D_DEBUG, "Enable debug output")
var timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

func signalHandler(signals chan os.Signal, done chan bool) {
	for _ = range signals {
		Log.Debug("[go-rtbh]: Sending cleanup signal to Amqp")
		Amqp.Control <- config.CTL_SHUTDOWN
		<-Amqp.Done

		Log.Debug("[go-rtbh]: Sending cleanup signal to Pipeline")
		Pipeline.Control <- config.CTL_SHUTDOWN
		<-Pipeline.Done

		done <- true
	}
}

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamps

	Log.Debug("Debug logging enabled")

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
	Log.Debug("Loaded " + strconv.Itoa(len(config.Ruleset)) + " rules")

	// Setup protocol library
	if err = proto.Setup(Log, Config); err != nil {
		Log.Fatal("[proto]: Initialization failed: " + err.Error())
	}
	Log.Debug("[proto]: Library initialized")

	// Setup redis client
	if Redis, err = proto.NewRedisClient(); err != nil {
		Log.Fatal("[Redis]: " + err.Error())
	}
	Log.Debug("[Redis]: Connected to " + Config.Redis.Address)

	// Setup lists library
	if err = lists.Setup(Log, Config, Redis); err != nil {
		Log.Fatal("[lists]: Initialization failed: " + err.Error())
	}
	Log.Debug("[lists]: Library initialized")

	// Configure pipeline
	if err = pipeline.Setup(Log, Config); err != nil {
		Log.Fatal("[pipeline]: Initialization failed: " + err.Error())
	}
	Log.Debug("[pipeline]: Library initialized")

	if Pipeline, err = pipeline.NewPipeline(); err != nil {
		Log.Fatal("[pipeline]: Failed to create new pipeline: " + err.Error())
	}
	Log.Debug("[pipeline]: Initialized new pipeline")

	// Configure AMQP client
	if Amqp, err = proto.NewAmqpClient(); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("[Amqp]: Connected to " + Config.Amqp.Address)
}

func main() {
	var input_data chan []byte

	input_data = make(chan []byte, config.D_INPUT_BUFSIZE)

	// Start signal handles
	signals = make(chan os.Signal, config.D_SIGNAL_BUFSIZE)
	allDone = make(chan bool)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, allDone)
	Log.Debug("[go-rtbh]: Started OS signal handler")

	// Start AMQP event slurper
	go Amqp.Slurp(input_data)
	Log.Debug("[go-rtbh]: Amqp event slurper started")

	// Start pipeline
	go Pipeline.Startup(input_data)
	Log.Debug("[go-rtbh]: Pipeline started")

	// Wait for program completion
	<-allDone
	Log.Debug("[go-rtbh]: Program finished")
}
