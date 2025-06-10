package main

import (
	"errors"
	"os/signal"
	"syscall"

	"codeberg.org/clambin/go-common/charmer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"
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

var arguments = charmer.Arguments{
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
