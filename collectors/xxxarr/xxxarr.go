package xxxarr

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient"
)

type xxxArrStats struct {
	version     string
	calendar    int
	queued      int
	monitored   int
	unmonitored int
}

type Updater struct {
	mediaclient.XXXArrAPI
}

func (updater *Updater) getStats() (interface{}, error) {
	var stats xxxArrStats
	var err error

	ctx := context.Background()

	stats.version, err = updater.GetVersion(ctx)

	if err == nil {
		stats.calendar, err = updater.GetCalendar(ctx)
	}

	if err == nil {
		stats.queued, err = updater.GetQueue(ctx)
	}

	if err == nil {
		stats.monitored, stats.unmonitored, err = updater.GetMonitored(ctx)
	}

	return stats, err
}
