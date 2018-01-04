package go_rtbhapp

import (
	"flag"
	"github.com/r3boot/go-rtbh/pkg/api"
	"github.com/r3boot/go-rtbh/pkg/blacklist"
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/history"
	"github.com/r3boot/go-rtbh/pkg/memcache"
	"github.com/r3boot/go-rtbh/pkg/orm"
	"github.com/r3boot/go-rtbh/pkg/whitelist"
	"github.com/r3boot/rlib/logger"
)

const MYNAME string = "go-rtbhapp"

var Config *config.Config
var Log logger.Log
var Blacklist *blacklist.Blacklist
var Whitelist *whitelist.Whitelist
var History *history.History
var ORM *orm.ORM
var API *api.RtbhApi

var cfgfile = flag.String("f", config.D_CFGFILE, "Configuration file to use")
var debug = flag.Bool("D", config.D_DEBUG, "Enable debug output")
var timestamps = flag.Bool("T", config.D_TIMESTAMP, "Enable timestamps in output")

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamps
	Log.Debug(MYNAME + ": Debug logging enabled")

	// Setup config library
	if err = config.Setup(Log); err != nil {
		Log.Fatal(err)
	}
	Config = config.New(*cfgfile)

	if err = events.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}

	if err = orm.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	ORM = orm.New()

	if err = memcache.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}

	if err = blacklist.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Blacklist = blacklist.New(nil)

	if err = whitelist.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	Whitelist = whitelist.New(nil)

	if err = history.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	History = history.New()

	if err = api.Setup(Log, Config); err != nil {
		Log.Fatal(err)
	}
	API = api.New(Blacklist, Whitelist, History)
	API.SetupRouting()
}

func main() {
	// Fire up the api
	API.RunServiceRoutine()
}
