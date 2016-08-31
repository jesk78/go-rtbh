package main

import (
	"errors"
	"flag"
	"github.com/r3boot/go-rtbh/lib/amqp"
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/listcache"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/pipeline"
	"github.com/r3boot/go-rtbh/lib/reaper"
	"github.com/r3boot/go-rtbh/lib/redis"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/lib/whitelist"
	"github.com/r3boot/rlib/logger"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const MYNAME string = "go-rtbh"

// Program libraries
var Config *config.Config
var Log logger.Log

var AmqpClient *amqp.AmqpClient
var RedisClient *redis.RedisClient
var Blacklist *blacklist.Blacklist
var Whitelist *whitelist.Whitelist
var History *history.History
var Resolver *resolver.Resolver
var Reaper *reaper.Reaper
var Pipeline *pipeline.Pipeline
var BGP bgp.BGP
var ORM *orm.ORM

// OS signals
var signals chan os.Signal
var allDone chan bool

// Command-line flags
var cfgfile = flag.String("f", config.D_CFGFILE, "Configuration file to use")
var debug = flag.Bool("D", config.D_DEBUG, "Enable debug output")
var timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

func signalHandler(signals chan os.Signal, done chan bool) {
	for _ = range signals {
		Log.Debug(MYNAME + ": Sending cleanup signal to Amqp")
		AmqpClient.Control <- config.CTL_SHUTDOWN
		<-AmqpClient.Done

		Log.Debug(MYNAME + ": Sending cleanup signal to Pipeline")
		Pipeline.Control <- config.CTL_SHUTDOWN
		<-Pipeline.Done

		Log.Debug(MYNAME + ": Sending cleanup signal to Reaper")
		Reaper.Control <- config.CTL_SHUTDOWN
		<-Reaper.Done

		if Config.General.Resolver.Enabled {
			Log.Debug(MYNAME + ": Sending cleanup signal to Resolver")
			Resolver.Control <- config.CTL_SHUTDOWN
			<-Resolver.Done
		}

		done <- true
	}
}

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamps
	Log.Debug(MYNAME + ": Debug logging enabled")

	// Setup random number generator
	rand.Seed(time.Now().UnixNano())

	// First, configure all dependencies
	if err = config.Setup(Log); err != nil {
		Log.Fatal(err)
	}
	Config = config.New(*cfgfile)

	if err = events.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}

	if Config.Redis.Address == "" && Config.Amqp.Address == "" {
		err = errors.New(MYNAME + ": No event feed to connect to")
		Log.Fatal(err)
	}

	if Config.Amqp.Address != "" {
		if err = amqp.Setup(Log, Config); err != nil {
			Log.Fatal(err)
		}
		AmqpClient = amqp.New()
	}

	if Config.Redis.Address != "" {
		if err = redis.Setup(Log, Config); err != nil {
			Log.Fatal(err)
		}
		RedisClient = redis.New()
	}

	if err = orm.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	ORM = orm.New()

	if err = bgp.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	BGP = bgp.New()

	// Then, setup all blacklist related libs
	if err = listcache.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}

	if err = blacklist.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Blacklist = blacklist.New(&BGP)

	if err = whitelist.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Whitelist = whitelist.New(&BGP)

	if err = history.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	History = history.New()

	if err = resolver.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	if Resolver, err = resolver.New(); err != nil {
		Log.Fatal(err)
	}

	if err = reaper.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Reaper = reaper.New(Blacklist)

	if err = pipeline.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Pipeline = pipeline.New(Blacklist, Whitelist, History)

	Log.Debug(MYNAME + ": All modules initialized")
}

func main() {
	var input_data chan []byte

	input_data = make(chan []byte, config.D_INPUT_BUFSIZE)

	// Start signal handles
	signals = make(chan os.Signal, config.D_SIGNAL_BUFSIZE)
	allDone = make(chan bool)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, allDone)
	Log.Debug(MYNAME + ": Started OS signal handler")

	// Start BGP routine
	go BGP.ServerRoutine()

	// Start pipeline
	go Pipeline.WorkManagerRoutine(input_data)

	// Start reaper
	go Reaper.CleanupExpiredRoutine()

	// Start DNS lookup thread
	go Resolver.UnknownLookupRoutine()

	// Start AMQP event slurper
	go AmqpClient.Slurp(input_data)

	// Wait for program completion
	<-allDone
	Log.Debug(MYNAME + ": Program finished")
}
