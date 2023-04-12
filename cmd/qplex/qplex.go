package main

import (
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/clambin/mediamon/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

var (
	cmd            *cobra.Command
	configFilename string
	tokenCmd       *cobra.Command
	viewsCmd       *cobra.Command
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", err)
		os.Exit(1)
	}
}

func init() {
	cmd = &cobra.Command{
		Use:     "qplex",
		Short:   "Plex utility",
		Version: version.BuildVersion,
	}
	cmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "qplex.yaml", "configuration file")
	cmd.PersistentFlags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "get an auth token",
		Run:   authToken,
	}
	cmd.AddCommand(tokenCmd)

	viewsCmd = &cobra.Command{
		Use:   "views",
		Short: "shows view counters for all media in all libraries",
		Run:   views,
	}
	cmd.AddCommand(viewsCmd)

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.AddConfigPath("/etc/qplex/")
		viper.AddConfigPath("$HOME/.qplex/")
		viper.AddConfigPath(".")
		viper.SetConfigName("qplex")
	}

	viper.SetDefault("debug", false)
	viper.SetDefault("authToken", "")
	viper.SetDefault("username", "")
	viper.SetDefault("password", "")
	viper.SetDefault("url", "")

	viper.SetEnvPrefix("QPLEX")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("failed to read config file", "err", err)
	}
}

func MakePlexClient() plex.Client {
	return plex.Client{
		HTTPClient: http.DefaultClient,
		URL:        viper.GetString("url"),
		AuthToken:  viper.GetString("auth.token"),
		UserName:   viper.GetString("auth.username"),
		Password:   viper.GetString("auth.password"),
		Product:    "qplex",
	}
}
