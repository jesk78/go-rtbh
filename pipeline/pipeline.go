package pipeline

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lib/blacklist"
	"github.com/r3boot/go-rtbh/lib/whitelist"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/rlib/logger"
	"regexp"
)

const MAX_WORKERS int = 8

type Pipeline struct {
	Control chan int
	Done    chan bool
}

var Log logger.Log
var Config *config.Config

var Whitelist *whitelist.Whitelist
var Blacklist *blacklist.Blacklist
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
	Whitelist = whitelist.New()
	Blacklist = blacklist.New()

	return
}

func (pl *Pipeline) Startup(input chan []byte) (err error) {
	var worker_queue chan chan []byte
	var stop_loop bool
	var worker_id int

	// Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

	worker_queue = make(chan chan []byte, MAX_WORKERS)

	// Startup event workers
	for worker_id = 1; worker_id <= MAX_WORKERS; worker_id++ {
		worker := NewEventWorker(worker_id, worker_queue)
		worker.Start()
	}

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case data := <-input:
			{
				go func() {
					worker := <-worker_queue
					worker <- data
				}()
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
