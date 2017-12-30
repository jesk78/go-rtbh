package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/r3boot/go-rtbh/lib/amqp"
	"github.com/r3boot/go-rtbh/lib/bgp"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/events"
	"github.com/r3boot/go-rtbh/lib/history"
	"github.com/r3boot/go-rtbh/lib/logger"
	"github.com/r3boot/go-rtbh/lib/orm"
	"github.com/r3boot/go-rtbh/lib/pipeline"
	"github.com/r3boot/go-rtbh/lib/reaper"
	"github.com/r3boot/go-rtbh/lib/redis"
	"github.com/r3boot/go-rtbh/lib/resolver"
	"github.com/r3boot/go-rtbh/lib/whitelist"
)

var (
	// Command-line flags
	cfgfile    = flag.String("f", config.D_CFGFILE, "Configuration file to use")
	debug      = flag.Bool("D", config.D_DEBUG, "Enable debug output")
	timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

	// Program libraries
	Config      *config.Config
	Logger      *logger.Logger
	AmqpClient  *amqp.AmqpClient
	RedisClient *redis.RedisClient
	Blacklist   *blacklist.Blacklist
	Whitelist   *whitelist.Whitelist
	History     *history.History
	Resolver    *resolver.Resolver
	Reaper      *reaper.Reaper
	Pipeline    *pipeline.Pipeline
	BGP         *bgp.BGP
	ORM         *orm.ORM

	// OS signals
	signals chan os.Signal
	allDone chan bool
)

func signalHandler(signals chan os.Signal, done chan bool) {
	for range signals {
		Logger.Debugf("main: Sending cleanup signal to Amqp")
		AmqpClient.Control <- config.CTL_SHUTDOWN
		<-AmqpClient.Done

		Logger.Debugf("main: Sending cleanup signal to Pipeline")
		Pipeline.Control <- config.CTL_SHUTDOWN
		<-Pipeline.Done

		Logger.Debugf("main: Sending cleanup signal to Reaper")
		Reaper.Control <- config.CTL_SHUTDOWN
		<-Reaper.Done

		if Config.General.Resolver.Enabled {
			Logger.Debugf("main: Sending cleanup signal to Resolver")
			Resolver.Control <- config.CTL_SHUTDOWN
			<-Resolver.Done
		}

		done <- true
	}
}

func init() {
	var err error

	// Setup random number generator
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	Logger = logger.NewLogger(*timestamps, *debug)
	Logger.Debugf("init: Debug logging enabled")

	// First, configure all dependencies
	Config, err := config.NewConfig(Logger, *cfgfile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	events.Setup(Logger, Config)

	if Config.Redis.Address == "" && Config.Amqp.Address == "" {
		Logger.Fatalf("init: No event feed to connect to")
	}

	if Config.Amqp.Address != "" {
		AmqpClient, err = amqp.NewAmqpClient(Logger, Config)
		if err != nil {
			Logger.Fatalf("init: %v", err)
		}
	}

	if Config.Redis.Address != "" {
		RedisClient = redis.NewRedisClient(Logger, Config)
	}

	ORM, err = orm.NewORM(Logger, Config)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	BGP, err = bgp.NewBGP(Logger, Config)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	// Then, setup all blacklist related libs
	Blacklist = blacklist.NewBlacklist(Logger, Config, BGP)
	Whitelist = whitelist.NewWhitelist(Logger, Config, BGP)
	History = history.NewHistory(Logger, Config)
	Pipeline = pipeline.NewPipeline(Logger, Config, Blacklist, Whitelist, History)

	Resolver, err = resolver.NewResolver(Logger, Config)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	Reaper, err = reaper.NewReaper(Logger, Config, Blacklist)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	Logger.Debugf("init: All modules initialized")
}

func main() {
	var input_data chan []byte

	input_data = make(chan []byte, config.D_INPUT_BUFSIZE)

	// Start signal handles
	signals = make(chan os.Signal, config.D_SIGNAL_BUFSIZE)
	allDone = make(chan bool)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, allDone)
	Logger.Debugf("main: Started OS signal handler")

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
	Logger.Debugf("main: Program finished")
}
