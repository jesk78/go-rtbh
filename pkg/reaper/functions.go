package reaper

import (
	"time"

	"github.com/r3boot/go-rtbh/pkg/config"
)

func (r *Reaper) CleanupExpiredRoutine() (err error) {
	t_now := time.Now()
	t_tick := t_now.Add(r.Interval)

	log.Debugf("Reaper: Running CleanupExpiredRoutine every %s", cfg.General.ReaperInterval)
	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case cmd := <-r.Control:
			{
				switch cmd {
				case config.CTL_SHUTDOWN:
					{
						log.Debugf("Reaper: Shutting down")
						stop_loop = true
						continue
					}
				}
			}
		default:
			{
				if t_now = time.Now(); t_now.After(t_tick) {
					r.blacklist.ReapExpiredEntries()
					t_tick = t_now.Add(r.Interval)
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}

	r.Done <- true

	return
}
