package pipeline

import (
	"fmt"
	"regexp"

	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/go-rtbh/lib/events"
)

type Worker struct {
	ID          int
	MyName      string
	parent      *Pipeline
	Ruleset     []*regexp.Regexp
	Work        chan []byte
	WorkerQueue chan chan []byte
	Done        chan bool
}

func NewEventWorker(id int, parent *Pipeline, workerQueue chan chan []byte) Worker {
	var worker Worker

	worker = Worker{
		ID:          id,
		MyName:      fmt.Sprintf("worker #%02d", id),
		parent:      parent,
		Ruleset:     config.Ruleset,
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

	log.Debugf("Worker.Start: starting worker %s", w.MyName)
	go func() {
		for {
			// Add this worker to the work queue
			w.WorkerQueue <- w.Work

			// Wait for new work
			select {
			case data = <-w.Work:
				{
					log.Debugf("Worker.%s: received new event", w.MyName)
					if event, err = events.NewEvent(data); err != nil {
						log.Warningf("Worker.%s: Failed to prepare event: %v", w.MyName, err)
						continue
					}

					if event.Address == "" {
						log.Warningf("Worker.%s: Failed to parse event %s", w.MyName, string(data))
						continue
					}

					if w.parent.whitelist.Listed(event.Address) {
						log.Warningf("Worker.%s: %s is on whitelist", w.MyName, event.Address)
						continue
					}

					if w.parent.blacklist.Listed(event.Address) {
						log.Warningf("Worker.%s: %s is already listed", w.MyName, event.Address)
						continue
					}

					if w.foundMatch(event.Reason) {
						event.ExpireIn = "5m"

						if err = w.parent.blacklist.Add(*event); err != nil {
							log.Warningf("Worker.%s: %v", w.MyName, err)
							continue
						}

						if err = w.parent.history.Add(*event); err != nil {
							log.Warningf("Worker.%s: %v", w.MyName, err)
							continue
						}

						log.Debugf("Worker.%s: Added %s to blacklist because of %s", w.MyName, event.Address, event.Reason)
					}
				}
			case <-w.Done:
				{
					return
				}
			}
		}
		log.Debugf("Worker.%s: Finished processing events", w.MyName)
	}()
}
