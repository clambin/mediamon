package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/clambin/mediamon/k8s/reaper"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	_ "k8s.io/client-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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

	slog.Info("restarter", "version", version.BuildVersion)
	go runPrometheusServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := check(ctx, *namespace, *name)

	if err != nil || *once {
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
			_ = check(ctx, *namespace, *name)
		}
	}
	ticker.Stop()
	slog.Info("scanner stopped")
}

func check(ctx context.Context, namespace, name string) error {
	r := reaper.Reaper{Connector: connect}
	deleted, err := r.Reap(ctx, namespace, name)
	if err == nil {
		deleteCounter.WithLabelValues(namespace, name).Add(float64(deleted))
	} else {
		slog.Error("scan failed", err)
	}
	return err
}

func connect() (kubernetes.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		// not running inside cluster. try to connect as external client
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("user home dir: %w", err)
		}
		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		slog.Debug("not running inside cluster. using kube config", "filename", kubeConfigPath)

		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("kubernetes config: %w", err)
		}
	}
	return kubernetes.NewForConfig(cfg)
}

func runPrometheusServer() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9091", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start prometheus metrics server", err)
	}
}
