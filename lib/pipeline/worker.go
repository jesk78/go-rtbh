package pipeline

import (
	"fmt"
	"github.com/r3boot/go-rtbh/events"
	"regexp"
)

type Worker struct {
	ID          int
	MyName      string
	Ruleset     []*regexp.Regexp
	Work        chan []byte
	WorkerQueue chan chan []byte
	Done        chan bool
}

func NewEventWorker(id int, ruleset []*regexp.Regexp, workerQueue chan chan []byte) Worker {
	var worker Worker

	worker = Worker{
		ID:          id,
		MyName:      fmt.Sprintf("worker #%02d", id),
		Ruleset:     ruleset,
		Work:        make(chan []byte),
		WorkerQueue: workerQueue,
		Done:        make(chan bool),
	}

	return worker
}

func (w *Worker) foundMatch(value string) bool {
	var re *regexp.Regexp
	for _, re = range w.Ruleset {
		if re.Match([]byte(value)) {
			return true
		}
	}

	return false
}

func (w *Worker) Start() {
	var data []byte
	var event *events.RTBHEvent
	var err error

	Log.Debug(MYNAME + "." + w.Myname + ": Starting up worker routine")
	go func() {
		for {
			// Add this worker to the work queue
			w.WorkerQueue <- w.Work

			// Wait for new work
			select {
			case data = <-w.Work:
				{
					Log.Debug(MYNAME + "." + w.Myname + ": Processing new event")
					if event, err = events.NewRTBHEvent(data); err != nil {
						Log.Warning(MYNAME + "." + w.Myname + ": Failed to prepare event" + err.Error())
						continue
					}

					if event.Address == "" {
						Log.Warning(MYNAME + "." + w.Myname + ": Failed to parse event: " + string(data))
						continue
					}

					if w.whitelist.Listed(event.Address) {
						Log.Warning(MYNAME + "." + w.Myname + ": Host " + event.Address + " is on whitelist")
						continue
					}

					if w.blacklist.Listed(event.Address) {
						Log.Warning(MYNAME + "." + w.Myname + ": Host " + event.Address + " is already listed")
						continue
					}

					if w.foundMatch(event.Reason) {
						event.ExpireIn = "1h"

						if err = w.blacklist.Add(*event); err != nil {
							Log.Warning(MYNAME + "." + w.Myname + ": Blacklist.Add failed: " + err.Error())
							continue
						}

						if err = w.history.Add(*event); err != nil {
							Log.Warning(MYNAME + "." + w.Myname + ": History.Add failed: " + err.Error())
							continue
						}

						Log.Debug(MYNAME + "." + w.MyName + ": Added " + event.Address + " to blacklist because of " + event.Reason)
					}
				}
			case <-w.Done:
				{
					return
				}
			}
		}
		Log.Debug(MYNAME + "." + w.MyName + ": Finished processing events")
	}()
}
