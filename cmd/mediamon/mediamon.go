package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
		"plex.username":                 {Default: ""},
		"plex.password":                 {Default: ""},
		"openvpn.connectivity.proxy":    {Default: ""},
		"openvpn.connectivity.interval": {Default: "10s"},
		"openvpn.bandwidth.filename":    {Default: ""},
	}
)

func main() {
	go func() {
		_ = http.ListenAndServe(":6060", nil)
	}()

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
	make func(string, string, *viper.Viper, *slog.Logger) (prometheus.Collector, error)
}

var constructors = map[string]constructor{
	"transmission.url": {
		name: "transmission",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return transmission.NewCollector(url, logger)
		},
	},
	"sonarr.url": {
		name: "sonarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewSonarrCollector(url, v.GetString("sonarr.apikey"), logger)
		},
	},
	"radarr.url": {
		name: "radarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewRadarrCollector(url, v.GetString("radarr.apikey"), logger)
		},
	},
	"prowlarr.url": {
		name: "prowlarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return prowlarr.New(url, v.GetString("prowlarr.apikey"), logger)
		},
	},
	"plex.url": {
		name: "plex",
		make: func(url, version string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return plex.NewCollector(
				version,
				url,
				v.GetString("plex.username"),
				v.GetString("plex.password"),
				logger,
			), nil
		},
	},
	"openvpn.connectivity.proxy": {
		name: "vpn connectivity",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			proxy, err := parseProxy(url)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy. connectivity won't be monitored: %w", err)
			}
			return connectivity.NewCollector(
				proxy,
				v.GetDuration("openvpn.connectivity.interval"),
				logger,
			), nil
		},
	},
	"openvpn.bandwidth.filename": {
		name: "vpn bandwidth",
		make: func(target, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return bandwidth.NewCollector(target, logger), nil
		},
	},
}

func createCollectors(version string, v *viper.Viper, logger *slog.Logger) []prometheus.Collector {
	collectors := make([]prometheus.Collector, 0, len(constructors))
	for key, c := range constructors {
		l := logger.With("collector", c.name)
		if value := v.GetString(key); value != "" {
			collector, err := c.make(value, version, v, l)
			if err != nil {
				l.Error("error creating collector", "err", err)
				continue
			}
			collectors = append(collectors, collector)
			l.Info("collector added", "source", value)
		}
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
