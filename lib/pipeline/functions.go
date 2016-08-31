package pipeline

import (
	"github.com/r3boot/go-rtbh/lib/config"
)

func (pl *Pipeline) WorkManagerRoutine(input chan []byte) (err error) {
	var worker_queue chan chan []byte
	var stop_loop bool
	var worker_id int

	// Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

	worker_queue = make(chan chan []byte, Config.General.NumWorkers)

	// Startup event workers
	for worker_id = 1; worker_id <= Config.General.NumWorkers; worker_id++ {
		worker := NewEventWorker(worker_id, pl, worker_queue)
		worker.Start()
	}

	Log.Debug(MYNAME + ": Starting WorkManagerRoutine")
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
