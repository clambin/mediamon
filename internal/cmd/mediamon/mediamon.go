package mediamon

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/go-common/charmer"
	"github.com/clambin/mediamon/v2/internal/collectors/bandwidth"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/clambin/mediamon/v2/internal/collectors/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFilename string
	RootCmd        = cobra.Command{
		Use:   "mediamon",
		Short: "Prometheus exporter for various media applications. Currently supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.",
		Run:   Main,
		PreRun: func(cmd *cobra.Command, args []string) {
			charmer.SetTextLogger(cmd, viper.GetBool("debug"))
		},
	}
)

func Main(cmd *cobra.Command, _ []string) {
	slog.SetDefault(charmer.GetLogger(cmd))

	slog.Info("mediamon starting", "version", cmd.Version)

	go func() {
		http.Handle(viper.GetString("metrics.path"), promhttp.Handler())
		if err := http.ListenAndServe(viper.GetString("metrics.addr"), nil); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start Prometheus listener", "err", err)
		}
	}()

	prometheus.MustRegister(createCollectors(cmd.Version)...)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer done()
	<-ctx.Done()

	slog.Info("mediamon exiting")
}

func createCollectors(version string) []prometheus.Collector {
	var collectors []prometheus.Collector

	if transmissionURL := viper.GetString("transmission.url"); transmissionURL != "" {
		slog.Info("monitoring Transmission", "url", transmissionURL)
		collectors = append(collectors, transmission.NewCollector(transmissionURL))
	}
	if sonarURL := viper.GetString("sonarr.url"); sonarURL != "" {
		slog.Info("monitoring Sonarr", "url", sonarURL)
		collectors = append(collectors, xxxarr.NewSonarrCollector(sonarURL, viper.GetString("sonarr.apikey")))
	}
	if radarURL := viper.GetString("radarr.url"); radarURL != "" {
		slog.Info("monitoring Radarr", "url", radarURL)
		collectors = append(collectors, xxxarr.NewRadarrCollector(radarURL, viper.GetString("radarr.apikey")))
	}
	if plexURL := viper.GetString("plex.url"); plexURL != "" {
		slog.Info("monitoring Plex", "url", plexURL)
		collectors = append(collectors, plex.NewCollector(
			version,
			plexURL,
			viper.GetString("plex.username"),
			viper.GetString("plex.password"),
		))
	}
	if proxyURL := viper.GetString("openvpn.connectivity.proxy"); proxyURL != "" {
		if proxy, err := parseProxy(proxyURL); err != nil {
			slog.Error("invalid proxy. connectivity won't be monitored", "err", err)
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
	RootCmd.Flags().StringVar(&configFilename, "config", "", "Configuration file")
	RootCmd.Flags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", RootCmd.Flags().Lookup("debug"))
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

	viper.SetDefault("debug", false)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.addr", ":9090")
	// viper.SetDefault("transmission.url", "http://transmission:9091/transmission/rpc")
	// viper.SetDefault("sonarr.url", "http://sonarr:8989")
	viper.SetDefault("sonarr.apikey", "")
	// viper.SetDefault("radarr.url", "http://radarr:7878")
	viper.SetDefault("radarr.apikey", "")
	// viper.SetDefault("plex.url", "http://plex:32400")
	viper.SetDefault("plex.username", "")
	viper.SetDefault("plex.password", "")
	// viper.SetDefault("openvpn.connectivity.proxy", "http://transmission:8888")
	viper.SetDefault("openvpn.connectivity.token", "")
	viper.SetDefault("openvpn.connectivity.interval", 5*time.Minute)
	// viper.SetDefault("openvpn.bandwidth.filename", "/data/client.status")

	viper.SetEnvPrefix("MEDIAMON")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "err", err)
		os.Exit(1)
	}
}
