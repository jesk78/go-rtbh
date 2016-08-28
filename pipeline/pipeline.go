package pipeline

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/events"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/rlib/logger"
	"regexp"
)

type Pipeline struct {
	Control chan int
	Done    chan bool
}

var Log logger.Log
var Config *config.Config

var Whitelist lists.Whitelist
var Blacklist *lists.Blacklist
var History lists.History
var Ruleset []*regexp.Regexp

func Setup(l logger.Log, cfg *config.Config) (err error) {
	Log = l
	Config = cfg

	return
}

func NewPipeline(ruleset []*regexp.Regexp) (pl *Pipeline, err error) {
	pl = &Pipeline{}
	Ruleset = ruleset

	return
}

func (pl *Pipeline) Startup(input chan []byte) (err error) {
	var stop_loop bool
	var event *events.RTBHEvent

	// Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

	Blacklist = lists.NewBlacklist()

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case data := <-input:
			{
				if event, err = events.NewRTBHEvent(data); err != nil {
					Log.Warning("[Pipeline] NewEvent: " + err.Error())
					continue
				}

				if event.Address == "" {
					// Log.Debug("[Pipeline]: Failed to parse event: " + string(data))
					continue
				}

				if Whitelist.Listed(event.Address) {
					Log.Warning("[Pipeline]: Host " + event.Address + " is on whitelist")
					continue
				}

				if Blacklist.Listed(event.Address) {
					Log.Warning("[Pipeline]: Host " + event.Address + " is already listed")
					History.Add(*event)
					continue
				}

				if FoundMatch(event.Reason) {
					event.ExpireIn = "1h"

					if err = Blacklist.Add(*event); err != nil {
						Log.Warning(err)
					}

					if err = History.Add(*event); err != nil {
						Log.Warning(err)
					}

					Log.Debug("[Pipeline]: Added " + event.Address + " to blacklist because of " + event.Reason)
					continue
				}
			}
		case cmd := <-pl.Control:
			{
				switch cmd {
				case config.CTL_SHUTDOWN:
					{
						Log.Debug("Shutting down pipeline")
						stop_loop = true
						continue
					}
				}
			}
		}
	}

	pl.Done <- true

	return
}
