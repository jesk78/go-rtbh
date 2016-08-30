package pipeline

import (
	"fmt"
	"github.com/r3boot/go-rtbh/events"
)

type Worker struct {
	ID          int
	MyName      string
	Work        chan []byte
	WorkerQueue chan chan []byte
	Done        chan bool
}

func NewEventWorker(id int, workerQueue chan chan []byte) Worker {
	var worker Worker

	worker = Worker{
		ID:          id,
		MyName:      fmt.Sprintf("worker #%02d", id),
		Work:        make(chan []byte),
		WorkerQueue: workerQueue,
		Done:        make(chan bool),
	}

	return worker
}

func (w *Worker) Debug(msg string) {
	Log.Debug("[" + w.MyName + "]: " + msg)
}

func (w *Worker) Warning(msg string) {
	Log.Warning("[" + w.MyName + "]: " + msg)
}

func (w *Worker) Start() {
	var data []byte
	var event *events.RTBHEvent
	var err error

	w.Debug("Starting up worker routine")
	go func() {
		for {
			// Add this worker to the work queue
			w.WorkerQueue <- w.Work

			// Wait for new work
			select {
			case data = <-w.Work:
				{
					w.Debug("Processing new event")
					if event, err = events.NewRTBHEvent(data); err != nil {
						w.Warning("Failed to prepare event" + err.Error())
						continue
					}

					if event.Address == "" {
						w.Warning("Failed to parse event: " + string(data))
						continue
					}

					if Whitelist.Listed(event.Address) {
						w.Warning("Host " + event.Address + " is on whitelist")
						continue
					}

					if Blacklist.Listed(event.Address) {
						w.Warning("Host " + event.Address + " is already listed")
						History.Add(*event)
						continue
					}

					if FoundMatch(event.Reason) {
						event.ExpireIn = "1h"

						if err = Blacklist.Add(*event); err != nil {
							w.Warning("Blacklist.Add failed: " + err.Error())
							continue
						}

						if err = History.Add(*event); err != nil {
							w.Warning("History.Add failed: " + err.Error())
							continue
						}

						w.Debug("Added " + event.Address + " to blacklist because of " + event.Reason)
					}
				}
			case <-w.Done:
				{
					return
				}
			}
		}
		w.Debug("Finished processing events")
	}()
}
