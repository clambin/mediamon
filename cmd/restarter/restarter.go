package main

import (
	"context"
	"errors"
	"flag"
	"github.com/clambin/mediamon/k8s/reaper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	_ "k8s.io/client-go"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	namespace = flag.String("namespace", "media", "namespace")
	name      = flag.String("name", "transmission", "deployment name")
	interval  = flag.Duration("interval", 5*time.Minute, "scanning interval")
	once      = flag.Bool("once", false, "scan once and exit")
	debug     = flag.Bool("debug", false, "enable debug mode")
)

var (
	deleteCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "restarter",
		Name:      "restarted",
		Help:      "total restarted pods",
	}, []string{"namespace", "name"})
)

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9091", nil)
		if !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("failed to start prometheus metrics server")
		}
	}()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := check(context.Background(), *namespace, *name); err != nil {
		log.WithError(err).Error("scan failed")
	}
	if *once {
		return
	}

	go scan(ctx, *interval)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

func scan(ctx context.Context, interval time.Duration) {
	log.WithField("interval", interval).Info("scanner started")
	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			if _, err := check(ctx, *namespace, *name); err != nil {
				log.WithError(err).Error("scan failed")
			}
		}
	}
	ticker.Stop()
	log.Infof("scanner stopped")
}

func check(ctx context.Context, namespace, name string) (int, error) {
	var r reaper.Reaper
	deleted, err := r.Reap(ctx, namespace, name)
	if err == nil {
		deleteCounter.WithLabelValues(namespace, name).Add(float64(deleted))
	}
	return deleted, err
}
