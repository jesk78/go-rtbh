package pipeline

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/go-rtbh/proto"
	"github.com/r3boot/rlib/logger"
	"regexp"
)

type Pipeline struct {
	Control chan int
	Done    chan bool
}

var Log logger.Log
var Config *config.Config

var Whitelist *lists.Whitelist
var Blacklist *lists.Blacklist
var History *lists.History
var Ruleset []*regexp.Regexp
var Bird *proto.Bird

func Setup(l logger.Log, cfg *config.Config) (err error) {
	Log = l
	Config = cfg

	return
}

func NewPipeline(ruleset []*regexp.Regexp) (pl *Pipeline, err error) {
	pl = &Pipeline{}
	Whitelist = lists.NewWhitelist()
	Blacklist = lists.NewBlacklist()
	Ruleset = ruleset
	Bird = proto.NewBirdClient()

	return
}

func (pl *Pipeline) Startup(input chan []byte) (err error) {
	var stop_loop bool
	var event *Event

	Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case data := <-input:
			{
				if event, err = NewEvent(data); err != nil {
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
					History.Update(event.Address)
					continue
				}

				if FoundMatch(event.Reason) {
					if !Blacklist.Add(event.Address, event.Reason) {
						continue
					}

					History.Update(event.Address)

					Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

					Log.Debug("[Pipeline]: Added " + event.Address + " to blacklist because of " + event.Reason)
					continue
				}

				Log.Debug(event)
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
