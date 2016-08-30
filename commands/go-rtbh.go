package go_rtbh

import (
	"flag"
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/orm"
	"github.com/r3boot/go-rtbh/pipeline"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/pg.v4"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// Program libraries
var Log logger.Log
var Config *config.Config
var Amqp *proto.AmqpClient
var Db *pg.DB
var Pipeline *pipeline.Pipeline
var Reaper *pipeline.Reaper
var Resolver *resolver.Resolver

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

		Log.Debug("[go-rtbh]: Sending cleanup signal to Reaper")
		Reaper.Control <- config.CTL_SHUTDOWN
		<-Reaper.Done

		Log.Debug("[go-rtbh]: Sending cleanup signal to Resolver")
		Resolver.Control <- config.CTL_SHUTDOWN
		<-Resolver.Done

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

	// Setup random number generator
	rand.Seed(time.Now().UnixNano())

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

	// Setup events
	events.Setup(Log)
	Log.Debug("[events]: Library initialized")

	// Setup protocol library
	if err = proto.Setup(Log, Config); err != nil {
		Log.Fatal("[proto]: Initialization failed: " + err.Error())
	}
	Log.Debug("[proto]: Library initialized")

	// Setup postgresql client
	if Db, err = proto.ConnectToPostgresql(); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("[postgresql]: Connected to " + Config.Database.Username + ":****@" + Config.Database.Address + "/" + Config.Database.Name)

	if err = orm.Setup(Log, Db); err != nil {
		Log.Fatal("[orm]: Initialization failed: " + err.Error())
	}
	Log.Debug("[orm]: Library initialized")

	// Setup lists library
	if err = lists.Setup(Log, Config); err != nil {
		Log.Fatal("[lists]: Initialization failed: " + err.Error())
	}
	Log.Debug("[lists]: Library initialized")

	// Configure pipeline
	if err = pipeline.Setup(Log, Config); err != nil {
		Log.Fatal("[pipeline]: Initialization failed: " + err.Error())
	}
	Log.Debug("[pipeline]: Library initialized")

	if Pipeline, err = pipeline.NewPipeline(config.Ruleset); err != nil {
		Log.Fatal("[pipeline]: Failed to create new pipeline: " + err.Error())
	}
	Log.Debug("[pipeline]: Initialized pipeline")

	if Reaper, err = pipeline.NewReaper("10s"); err != nil {
		Log.Fatal("[pipeline]: Failed to create new reaper: " + err.Error())
	}
	Log.Debug("[pipeline]: Initialized reaper")

	if Resolver, err = resolver.NewResolver(); err != nil {
		Log.Fatal("[Resolver]: Failed to create dns resolver: " + err.Error())
	}

	// Configure AMQP client
	if Amqp, err = proto.NewAmqpClient(); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("[Amqp]: Connected to " + Config.Amqp.Address)

	// Configure BGP routine
	if err = proto.ConfigureBGP(); err != nil {
		Log.Fatal(err)
	}
	Log.Debug("[BGP]: Library initialized")
}

func main() {
	var bgpPeer config.BGPPeer
	var input_data chan []byte

	input_data = make(chan []byte, config.D_INPUT_BUFSIZE)

	// Start signal handles
	signals = make(chan os.Signal, config.D_SIGNAL_BUFSIZE)
	allDone = make(chan bool)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, allDone)
	Log.Debug("[go-rtbh]: Started OS signal handler")

	// Start BGP routine
	proto.RunBGP()
	Log.Debug("[go-rtbh]: BGP route injector started")

	// Add BGP neighbors
	for _, bgpPeer = range Config.BGP.Peers {
		Log.Debug("[go-rtbh] Adding BGP neighbor " + bgpPeer.Address + " as " + Config.BGP.Asnum)
		proto.AddBGPNeighbor(bgpPeer.Address)
	}

	// Start pipeline
	go Pipeline.Startup(input_data)
	Log.Debug("[go-rtbh]: Event pipeline started")

	// Start reaper
	go Reaper.Startup()
	Log.Debug("[go-rtbh]: Blacklist reaper started")

	// Start DNS lookup thread
	go Resolver.UnknownLookupRoutine()
	Log.Debug("[go-rtbh]: Dns lookup routine started")

	// Start AMQP event slurper
	go Amqp.Slurp(input_data)
	Log.Debug("[go-rtbh]: Amqp event slurper started")

	// Wait for program completion
	<-allDone
	Log.Debug("[go-rtbh]: Program finished")
}
