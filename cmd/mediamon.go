package main

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/clambin/go-common/charmer"
	"github.com/clambin/mediamon/v2/internal/collectors/bandwidth"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/clambin/mediamon/v2/internal/collectors/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/prowlarr"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version        = "change_me"
	configFilename string
	rootCmd        = cobra.Command{
		Use:   "mediamon",
		Short: "Prometheus exporter for various media applications. Currently supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.",
		Run:   Main,
		PreRun: func(cmd *cobra.Command, args []string) {
			charmer.SetTextLogger(cmd, viper.GetBool("debug"))
		},
	}

	arguments = charmer.Arguments{
		"debug":                         {Default: false},
		"metrics.path":                  {Default: "/metrics"},
		"metrics.addr":                  {Default: ":9090"},
		"transmission.url":              {Default: ""},
		"sonarr.url":                    {Default: ""},
		"sonarr.apikey":                 {Default: ""},
		"radarr.url":                    {Default: ""},
		"radarr.apikey":                 {Default: ""},
		"plex.url":                      {Default: ""},
		"plex.client-id":                {Default: ""},
		"plex.username":                 {Default: ""},
		"plex.password":                 {Default: ""},
		"plex.jwt.enable":               {Default: false},
		"plex.jwt.path":                 {Default: ""},
		"plex.jwt.passphrase":           {Default: ""},
		"openvpn.connectivity.proxy":    {Default: ""},
		"openvpn.connectivity.interval": {Default: "10s"},
		"openvpn.bandwidth.filename":    {Default: ""},
	}
)

func main() {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func Main(cmd *cobra.Command, _ []string) {
	logger := charmer.GetLogger(cmd)
	slog.SetDefault(logger)

	logger.Info("mediamon starting", "version", cmd.Version, "addr", viper.GetString("metrics.addr"))

	go func() {
		http.Handle(viper.GetString("metrics.path"), promhttp.Handler())
		if err := http.ListenAndServe(viper.GetString("metrics.addr"), nil); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start Prometheus listener", "err", err)
		}
	}()

	prometheus.MustRegister(createCollectors(cmd.Version, viper.GetViper(), logger)...)

	ctx, done := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer done()
	<-ctx.Done()

	logger.Info("mediamon exiting")
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&configFilename, "config", "", "Configuration file")
	_ = charmer.SetPersistentFlags(&rootCmd, viper.GetViper(), arguments)
	_ = charmer.SetDefaults(viper.GetViper(), arguments)
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.AddConfigPath("/etc/mediamon/")
		viper.AddConfigPath("$HOME/.mediamon")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("MEDIAMON")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
	}
}

type constructor struct {
	name string
}

var constructors = map[string]constructor{
	"transmission.url": {
		name: "transmission",
	},
	"sonarr.url": {
		name: "sonarr",
	},
	"radarr.url": {
		name: "radarr",
	},
	"prowlarr.url": {
		name: "prowlarr",
	},
	"plex.url": {
		name: "plex",
	},
	"openvpn.connectivity.proxy": {
		name: "connectivity",
	},
	"openvpn.bandwidth.filename": {
		name: "bandwidth",
	},
}

func createCollectors(version string, v *viper.Viper, logger *slog.Logger) []prometheus.Collector {
	collectors := make([]prometheus.Collector, 0, len(constructors))
	for key, c := range constructors {
		target := v.GetString(key)
		if target == "" {
			continue
		}
		l := logger.With("collector", c.name)

		// for connectivity, we need a proxy-enabled http.Transport.
		var rt http.RoundTripper
		if key == "openvpn.connectivity.proxy" {
			proxy, err := parseProxy(target)
			if err != nil {
				l.Error("failed to parse proxy URL. connectivity monitoring disabled", "err", err)
				continue
			}
			rt = &http.Transport{Proxy: http.ProxyURL(proxy)}
		}

		var collector prometheus.Collector
		var err error
		httpClient, metrics := instrumentedHttpClient(c.name, rt)
		collectors = append(collectors, metrics)

		switch key {
		case "transmission.url":
			collector, err = transmission.NewCollector(httpClient, target, l)
		case "sonarr.url":
			collector, err = xxxarr.NewSonarrCollector(target, v.GetString("sonarr.apikey"), httpClient, l)
		case "radarr.url":
			collector, err = xxxarr.NewRadarrCollector(target, v.GetString("radarr.apikey"), httpClient, l)
		case "prowlarr.url":
			collector, err = prowlarr.New(target, v.GetString("prowlarr.apikey"), httpClient, l)
		case "plex.url":
			pcfg := plex.Config{
				UserName:      v.GetString("plex.username"),
				Password:      v.GetString("plex.password"),
				ClientID:      v.GetString("plex.client-id"),
				UseJWT:        v.GetBool("plex.jwt.enable"),
				JWTLocation:   v.GetString("plex.jwt.path"),
				JWTPassphrase: v.GetString("plex.jwt.passphrase"),
				Version:       version,
			}
			collector = plex.NewCollector(target, pcfg, httpClient, l)
		case "openvpn.bandwidth.filename":
			collector = bandwidth.NewCollector(target, l)
		case "openvpn.connectivity.proxy":
			collector = connectivity.NewCollector(httpClient, v.GetDuration("openvpn.connectivity.interval"), l)
		}
		if err != nil {
			l.Error("error creating collector", "err", err)
			continue
		}
		collectors = append(collectors, collector)
		l.Info("collector added", "source", target)
	}
	return collectors
}

func parseProxy(proxyURL string) (*url.URL, error) {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	if proxy.Scheme == "" || proxy.Host == "" {
		return nil, fmt.Errorf("missing scheme / host")
	}
	return proxy, nil
}

func instrumentedHttpClient(application string, roundTripper http.RoundTripper) (*http.Client, prometheus.Collector) {
	metrics := requestMetrics{
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   "mediamon",
			Subsystem:   "http",
			Name:        "requests_total",
			Help:        "total number of http requests",
			ConstLabels: prometheus.Labels{"application": application},
		}, []string{"method", "code"}),
		latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:   "mediamon",
			Subsystem:   "http",
			Name:        "request_duration_seconds",
			Help:        "duration of http requests",
			ConstLabels: prometheus.Labels{"application": application},
		}, []string{"method", "code"}),
	}

	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: promhttp.InstrumentRoundTripperCounter(metrics.counter,
			promhttp.InstrumentRoundTripperDuration(metrics.latency,
				cmp.Or(roundTripper, http.DefaultTransport),
			),
		),
	}
	return &client, &metrics
}

var _ prometheus.Collector = &requestMetrics{}

type requestMetrics struct {
	counter *prometheus.CounterVec
	latency *prometheus.SummaryVec
}

func (r *requestMetrics) Describe(ch chan<- *prometheus.Desc) {
	r.counter.Describe(ch)
	r.latency.Describe(ch)
}

func (r *requestMetrics) Collect(ch chan<- prometheus.Metric) {
	r.counter.Collect(ch)
	r.latency.Collect(ch)
}
