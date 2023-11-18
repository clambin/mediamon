package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/internal/qplex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	rootCmd        *cobra.Command
	configFilename string
)

var BuildVersion = "change_me"

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd = &cobra.Command{
		Use:     "qplex",
		Short:   "Plex utility",
		Version: BuildVersion,
	}
	rootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "qplex.yaml", "configuration file")
	rootCmd.PersistentFlags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "get an auth token",
		Run:   getAuthToken,
	}
	rootCmd.AddCommand(tokenCmd)

	viewsCmd := &cobra.Command{
		Use:   "views",
		Short: "list view counters for all media in all libraries",
		Run:   getViews,
	}
	viewsCmd.Flags().BoolP("reverse", "r", false, "Sort view count high to low")
	_ = viper.BindPFlag("views.reverse", viewsCmd.Flags().Lookup("reverse"))
	viewsCmd.Flags().BoolP("server", "s", false, "Use server token to query all users")
	_ = viper.BindPFlag("views.server", viewsCmd.Flags().Lookup("server"))
	rootCmd.AddCommand(viewsCmd)

	sessionCmd := &cobra.Command{
		Use:   "sessions",
		Short: "list active sessions",
		Run:   getSessions,
	}
	rootCmd.AddCommand(sessionCmd)

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

func getAuthToken(cmd *cobra.Command, _ []string) {
	c := plex.New(
		viper.GetString("auth.username"),
		viper.GetString("auth.password"),
		"qplex",
		cmd.Version,
		viper.GetString("url"),
		nil,
	)
	if authToken := viper.GetString("auth.token"); authToken != "" {
		c.SetAuthToken(authToken)
	}

	token, err := c.GetAuthToken(context.Background())
	if err != nil {
		slog.Error("failed to get authentication token", "err", err)
		return
	}
	fmt.Printf("authToken: %s\n", token)
}

func getViews(cmd *cobra.Command, _ []string) {
	ctx := context.Background()
	tokens, err := getTokens(ctx, viper.GetBool("views.server"))
	if err != nil {
		slog.Error("failed to get tokens", "err", err)
		return
	}
	c := plex.New(
		"",
		"",
		"qplex",
		cmd.Version,
		viper.GetString("url"),
		nil,
	)

	views, err := qplex.GetViews(ctx, c, tokens, viper.GetBool("views.reverse"))
	if err != nil {
		slog.Error("failed to get views", "err", err)
		return
	}

	if len(views) > 0 {
		fmt.Printf("%-20s %-60s %s\n", "LIBRARY", "TITLE", "VIEWS")
		for _, entry := range views {
			fmt.Printf("%-20s %-60s %d\n", entry.Library, entry.Title, entry.Views)
		}
	}
}

func getTokens(ctx context.Context, server bool) ([]string, error) {
	if server {
		return getServerTokens(ctx)
	}

	c := plex.New(
		viper.GetString("auth.username"),
		viper.GetString("auth.password"),
		"qplex",
		BuildVersion,
		viper.GetString("url"),
		nil,
	)
	if authToken := viper.GetString("auth.token"); authToken != "" {
		c.SetAuthToken(authToken)
	}

	token, err := c.GetAuthToken(ctx)
	return []string{token}, err
}

func getServerTokens(ctx context.Context) ([]string, error) {
	serverToken := viper.GetString("auth.serverToken")
	if serverToken == "" {
		return nil, fmt.Errorf("no server token configured")
	}

	tokens, err := GetAccessTokens(ctx, serverToken)
	if err != nil {
		return nil, err
	}

	serverTokens := []string{serverToken}
	for _, token := range tokens {
		if token.Type == "server" {
			slog.Debug("token found", "user", token.Invited.Title)
			serverTokens = append(serverTokens, token.Token)
		}
	}

	return serverTokens, nil
}

func GetAccessTokens(ctx context.Context, serverToken string) ([]AccessToken, error) {
	args := make(url.Values)
	args.Set("auth_token", serverToken)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://plex.tv/api/v2/server/access_tokens?"+args.Encode(), nil)
	req.Header.Set("Accept", "application/json")

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var tokens []AccessToken
	err = json.NewDecoder(resp.Body).Decode(&tokens)
	return tokens, err
}

type AccessToken struct {
	Type      string    `json:"type"`
	Device    string    `json:"device,omitempty"`
	Token     string    `json:"token"`
	Owned     bool      `json:"owned"`
	CreatedAt time.Time `json:"createdAt"`
	Invited   struct {
		Id       int         `json:"id"`
		Uuid     string      `json:"uuid"`
		Title    string      `json:"title"`
		Username interface{} `json:"username"`
		Thumb    string      `json:"thumb"`
		Profile  struct {
			AutoSelectAudio              bool        `json:"autoSelectAudio"`
			DefaultAudioLanguage         interface{} `json:"defaultAudioLanguage"`
			DefaultSubtitleLanguage      interface{} `json:"defaultSubtitleLanguage"`
			AutoSelectSubtitle           int         `json:"autoSelectSubtitle"`
			DefaultSubtitleAccessibility int         `json:"defaultSubtitleAccessibility"`
			DefaultSubtitleForced        int         `json:"defaultSubtitleForced"`
		} `json:"profile"`
		Scrobbling    []interface{} `json:"scrobbling"`
		ScrobbleTypes string        `json:"scrobbleTypes"`
	} `json:"invited,omitempty"`
	Settings struct {
		AllowChannels      bool        `json:"allowChannels"`
		FilterMovies       *string     `json:"filterMovies"`
		FilterMusic        *string     `json:"filterMusic"`
		FilterPhotos       interface{} `json:"filterPhotos"`
		FilterTelevision   *string     `json:"filterTelevision"`
		FilterAll          interface{} `json:"filterAll"`
		AllowSync          bool        `json:"allowSync"`
		AllowCameraUpload  bool        `json:"allowCameraUpload"`
		AllowSubtitleAdmin bool        `json:"allowSubtitleAdmin"`
		AllowTuners        int         `json:"allowTuners"`
	} `json:"settings,omitempty"`
	Sections []struct {
		Key       int       `json:"key"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"sections,omitempty"`
}

func getSessions(cmd *cobra.Command, _ []string) {
	c := plex.New(
		viper.GetString("auth.username"),
		viper.GetString("auth.password"),
		"qplex",
		cmd.Version,
		viper.GetString("url"),
		nil,
	)
	if authToken := viper.GetString("auth.token"); authToken != "" {
		c.SetAuthToken(authToken)
	}

	sessions, err := c.GetSessions(context.Background())
	if err != nil {
		slog.Error("failed to get active sessions", "err", err)
		return
	}

	if len(sessions) > 0 {
		fmt.Printf("%-10s %-40s %-8s %-10s %-10s %s\n", "USER", "TITLE", "LOCATION", "VIDEO MODE", "STATE", "PROGRESS")
	}
	for _, session := range sessions {
		video := session.TranscodeSession.VideoDecision
		if video == "" {
			video = "direct"
		}
		fmt.Printf("%-10s %-40s %-8s %-10s %-10s %.2f%%\n",
			session.User.Title,
			session.GetTitle(),
			session.Session.Location,
			video,
			session.Player.State,
			100.0*session.GetProgress(),
		)
	}
}
