package reaper

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"time"
)

func (r *Reaper) CleanupExpiredRoutine() (err error) {
	var stop_loop bool
	var t_tick time.Time
	var t_now time.Time

	t_now = time.Now()
	t_tick = t_now.Add(r.Interval)

	stop_loop = false
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
						Log.Debug("Shutting down reaper")
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
