package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/version"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// API interface
type API interface {
	GetVersion(context.Context) (string, error)
	GetSessions(ctx context.Context) (sessions []Session, err error)
}

// Client calls the Plex APIs
type Client struct {
	Client    *http.Client
	URL       string
	Options   Options
	AuthURL   string
	UserName  string
	Password  string
	authToken string
}

// Options contains options to alter Client behaviour
type Options struct {
	PrometheusMetrics metrics.APIClientMetrics
}

// GetVersion retrieves the version of the Plex server
func (client *Client) GetVersion(ctx context.Context) (version string, err error) {
	var stats identityStats
	err = client.call(ctx, "/identity", &stats)
	if err == nil {
		version = stats.MediaContainer.Version
	}
	return
}

// Session represents a user watching a stream on Plex
type Session struct {
	Title     string  // title of the movie, tv show
	User      string  // Name of user
	Local     bool    // Is the user local (LAN) or not (WAN)?
	Transcode bool    // Does the session transcode the video?
	Throttled bool    // Is transcoding currently throttled?
	Speed     float64 // Current transcoding speed
}

// GetSessions retrieves session information from the server.
func (client *Client) GetSessions(ctx context.Context) (sessions []Session, err error) {
	var stats sessionStats
	err = client.call(ctx, "/status/sessions", &stats)
	if err == nil {
		for _, entry := range stats.MediaContainer.Metadata {
			sessions = append(sessions, Session{
				Title:     entry.GrandparentTitle,
				User:      entry.User.Title,
				Local:     entry.Player.Local,
				Transcode: entry.TranscodeSession.VideoDecision == "transcode",
				Throttled: entry.TranscodeSession.Throttled,
				Speed:     entry.TranscodeSession.Speed,
			})
		}
	}
	return
}

// call the specified Plex API endpoint
func (client *Client) call(ctx context.Context, endpoint string, response interface{}) (err error) {
	defer func() {
		client.Options.PrometheusMetrics.ReportErrors(err, "plex", endpoint)
	}()

	client.authToken, err = client.authenticate(ctx)

	if err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Plex-Token", client.authToken)

	timer := client.Options.PrometheusMetrics.MakeLatencyTimer("plex", endpoint)

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	return
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (client *Client) authenticate(ctx context.Context) (authToken string, err error) {
	defer func() {
		client.Options.PrometheusMetrics.ReportErrors(err, "plex", "auth")
	}()

	if authToken = client.authToken; authToken != "" {
		return
	}

	authBody := fmt.Sprintf("user%%5Blogin%%5D=%s&user%%5Bpassword%%5D=%s",
		client.UserName,
		client.Password,
	)

	authURL := "https://plex.tv/users/sign_in.xml"
	if client.AuthURL != "" {
		authURL = client.AuthURL
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, authURL, bytes.NewBufferString(authBody))
	req.Header.Add("X-Plex-Product", "github.com/clambin/mediamon")
	req.Header.Add("X-Plex-Version", version.BuildVersion)
	req.Header.Add("X-Plex-Client-Identifier", "github.com/clambin/mediamon-v"+version.BuildVersion)

	timer := client.Options.PrometheusMetrics.MakeLatencyTimer("plex", "auth")

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("%s", resp.Status)
	}

	// TODO: there's three different places in the response where the authToken appears.
	// Which is the officially supported version?
	var authResponse struct {
		XMLName             xml.Name `xml:"user"`
		AuthenticationToken string   `xml:"authenticationToken,attr"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err = xml.Unmarshal(body, &authResponse); err == nil {
		authToken = authResponse.AuthenticationToken
	}

	log.WithFields(log.Fields{
		"err":       err,
		"authToken": authToken,
	}).Debug("plex authenticate")

	return
}
