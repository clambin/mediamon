package scheduler

import (
	"time"
)

type Scheduler struct {
	Schedule chan *ScheduledTask
	Stop     chan struct{}

	scheduled []Scheduled
}

type ScheduledTask struct {
	Task     Runnable
	Interval time.Duration
}

func New() *Scheduler {
	return &Scheduler{
		Schedule:  make(chan *ScheduledTask),
		Stop:      make(chan struct{}),
		scheduled: make([]Scheduled, 0),
	}
}

func (scheduler *Scheduler) schedule(scheduledTask *ScheduledTask) {
	scheduled := Scheduled{
		stop:   make(chan struct{}),
		task:   scheduledTask.Task,
		ticker: time.NewTicker(scheduledTask.Interval),
	}
	go scheduled.Run()
	scheduler.scheduled = append(scheduler.scheduled, scheduled)
}

func (scheduler *Scheduler) Run() {
loop:
	for {
		select {
		case <-scheduler.Stop:
			break loop
		case scheduledTask := <-scheduler.Schedule:
			scheduler.schedule(scheduledTask)
		}
	}

	for _, scheduledTask := range scheduler.scheduled {
		scheduledTask.stop <- struct{}{}
	}
}
