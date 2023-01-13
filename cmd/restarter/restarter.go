package main

import (
	"context"
	"errors"
	"flag"
	"github.com/clambin/mediamon/k8s/reaper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
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

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stdout)))

	go runPrometheusServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := check(ctx, *namespace, *name); err != nil {
		slog.Error("scan failed", err)
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
	slog.Info("scanner started", "interval", interval)
	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			if _, err := check(ctx, *namespace, *name); err != nil {
				slog.Error("scan failed", err)
			}
		}
	}
	ticker.Stop()
	slog.Info("scanner stopped")
}

func check(ctx context.Context, namespace, name string) (int, error) {
	r := reaper.Reaper{}
	deleted, err := r.Reap(ctx, namespace, name)
	if err == nil {
		deleteCounter.WithLabelValues(namespace, name).Add(float64(deleted))
	}
	return deleted, err
}

func runPrometheusServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9091", nil)
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start prometheus metrics server", err)
			panic(err)
		}
	}()
}
