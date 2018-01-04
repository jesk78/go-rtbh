package pipeline

import (
	"github.com/r3boot/go-rtbh/pkg/config"
)

func (pl *Pipeline) WorkManagerRoutine(input chan []byte) (err error) {
	// Bird.ExportPrefixes(Whitelist.GetAll(), Blacklist.GetAll())

	worker_queue := make(chan chan []byte, cfg.General.NumWorkers)

	// Startup event workers
	for worker_id := 1; worker_id <= cfg.General.NumWorkers; worker_id++ {
		worker := NewEventWorker(worker_id, pl, worker_queue)
		worker.Start()
	}

	log.Debugf("Pipeline.WorkManagerRoutine: Starting pipeline")
	stop_loop := false
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
						log.Debugf("Pipeline.WorkManagerRoutine: Shutting down pipeline")
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
