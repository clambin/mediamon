package scheduler

import "time"

type Runnable interface {
	Run() error
}

type Scheduled struct {
	stop   chan struct{}
	task   Runnable
	ticker *time.Ticker
}

func (scheduled *Scheduled) Run() {
	_ = scheduled.task.Run()
loop:
	for {
		select {
		case <-scheduled.ticker.C:
			_ = scheduled.task.Run()
		case <-scheduled.stop:
			break loop
		}
	}
}
