package scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	Schedule chan *ScheduledTask

	scheduled []Scheduled
}

type ScheduledTask struct {
	Task     Runnable
	Interval time.Duration
}

func New() *Scheduler {
	return &Scheduler{
		Schedule:  make(chan *ScheduledTask),
		scheduled: make([]Scheduled, 0),
	}
}

func (scheduler *Scheduler) schedule(ctx context.Context, scheduledTask *ScheduledTask) {
	scheduled := Scheduled{
		task:   scheduledTask.Task,
		ticker: time.NewTicker(scheduledTask.Interval),
	}
	go scheduled.Run(ctx)
	scheduler.scheduled = append(scheduler.scheduled, scheduled)
}

func (scheduler *Scheduler) Run(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case scheduledTask := <-scheduler.Schedule:
			scheduler.schedule(ctx, scheduledTask)
		}
	}
}
