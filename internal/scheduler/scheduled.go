package scheduler

import (
	"context"
	"time"
)

type Runnable interface {
	Run(ctx context.Context) error
}

type Scheduled struct {
	task   Runnable
	ticker *time.Ticker
}

func (scheduled *Scheduled) Run(ctx context.Context) {
	_ = scheduled.task.Run(ctx)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-scheduled.ticker.C:
			_ = scheduled.task.Run(ctx)
		}
	}
}
