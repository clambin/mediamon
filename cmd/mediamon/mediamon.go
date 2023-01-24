package main

import (
	"errors"
	"fmt"
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFilename string
	cmd            = cobra.Command{
		Use:   "mediamon",
		Short: "Prometheus exporter for various media applications. Currently supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.",
		Run:   Main,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", err)
		os.Exit(1)
	}
}

func Main(_ *cobra.Command, _ []string) {
	var opts slog.HandlerOptions
	if viper.GetBool("debug") {
		opts.Level = slog.LevelDebug
		//opts.AddSource = true
	}
	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stderr)))

	slog.Info("mediamon starting", "version", version.BuildVersion)

	go runPrometheusServer()

	collectors := createCollectors()
	prometheus.MustRegister(collectors...)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	slog.Info("mediamon exiting")
}

func runPrometheusServer() {
	http.Handle(viper.GetString("metrics.path"), promhttp.Handler())
	if err := http.ListenAndServe(viper.GetString("metrics.addr"), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", err)
		panic(err)
	}
}

func createCollectors() []prometheus.Collector {
	var collectors []prometheus.Collector

	if url := viper.GetString("transmission.url"); url != "" {
		slog.Info("monitoring Transmission", "url", url)
		collectors = append(collectors, transmission.NewCollector(url))
	}
	if url := viper.GetString("sonarr.url"); url != "" {
		slog.Info("monitoring Sonarr", "url", url)
		collectors = append(collectors, xxxarr.NewSonarrCollector(url, viper.GetString("sonarr.apikey")))
	}
	if url := viper.GetString("radarr.url"); url != "" {
		slog.Info("monitoring Radarr", "url", url)
		collectors = append(collectors, xxxarr.NewRadarrCollector(url, viper.GetString("radarr.apikey")))
	}
	if url := viper.GetString("plex.url"); url != "" {
		slog.Info("monitoring Plex", "url", url)
		collectors = append(collectors, plex.NewCollector(url,
			viper.GetString("plex.username"),
			viper.GetString("plex.password"),
		))
	}
	if proxyUrl := viper.GetString("openvpn.connectivity.proxy"); proxyUrl != "" {
		proxy, err := parseProxy(proxyUrl)
		if err != nil {
			slog.Error("invalid proxy. connectivity won't be monitored", err)
		} else {
			slog.Info("monitoring openVPN connectivity", "url", proxy.String())
			collectors = append(collectors, connectivity.NewCollector(
				viper.GetString("openvpn.connectivity.token"),
				proxy,
				viper.GetDuration("openvpn.connectivity.interval"),
			))
		}
	}
	if filename := viper.GetString("openvpn.bandwidth.filename"); filename != "" {
		slog.Info("monitoring openVPN bandwidth", "filename", filename)
		collectors = append(collectors, bandwidth.NewCollector(filename))
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

func init() {
	cobra.OnInitialize(initConfig)
	cmd.Version = version.BuildVersion
	cmd.Flags().StringVar(&configFilename, "config", "", "Configuration file")
	cmd.Flags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
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

	// TODO: we probably don't want to default the urls's: otherwise all services are on by default and switching them off becomes cumbersome (or impossible?)
	viper.SetDefault("debug", false)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.addr", ":9090")
	viper.SetDefault("transmission.url", "http://transmission:9091/transmission/rpc")
	viper.SetDefault("sonarr.url", "http://sonarr:8989")
	viper.SetDefault("sonarr.apikey", "")
	viper.SetDefault("radarr.url", "http://radarr:7878")
	viper.SetDefault("radarr.apikey", "")
	viper.SetDefault("plex.url", "http://plex:32400")
	viper.SetDefault("plex.username", "")
	viper.SetDefault("plex.password", "")
	viper.SetDefault("openvpn.connectivity.proxy", "http://transmission:8888")
	viper.SetDefault("openvpn.connectivity.token", "")
	viper.SetDefault("openvpn.connectivity.interval", 5*time.Minute)
	viper.SetDefault("openvpn.bandwidth.filename", "/data/client.status")

	viper.SetEnvPrefix("MEDIAMON")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", err)
		os.Exit(1)
	}
}
