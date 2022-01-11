package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clambin/mediamon/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	PrometheusSummary *prometheus.SummaryVec
}

// GetVersion retrieves the version of the Plex server
func (client *Client) GetVersion(ctx context.Context) (version string, err error) {
	var resp []byte
	if resp, err = client.call(ctx, "/identity"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))

		var stats struct {
			MediaContainer struct {
				Version string
			}
		}
		err = decoder.Decode(&stats)

		if err == nil {
			version = stats.MediaContainer.Version
		}
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": version,
	}).Debug("plex getVersion")

	return
}

type Session struct {
	Title     string
	User      string
	Local     bool
	Transcode bool
	Throttled bool
	Speed     float64
}

type sessionStats struct {
	MediaContainer struct {
		Metadata []struct {
			GrandparentTitle string
			Media            []struct {
				Part []struct {
					Stream []struct {
						Decision string
						Location string
					}
				}
			}
			User struct {
				Title string
			}
			Player struct {
				Local bool
			}
			TranscodeSession struct {
				Throttled     bool
				Speed         float64
				VideoDecision string
			}
		}
	}
}

// GetSessions retrieves session information from the server.
func (client *Client) GetSessions(ctx context.Context) (sessions []Session, err error) {
	var resp []byte
	if resp, err = client.call(ctx, "/status/sessions"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))

		var stats sessionStats
		err = decoder.Decode(&stats)

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

	log.WithError(err).Debug("plex getSessions")

	return
}

// call the specified Plex API endpoint
func (client *Client) call(ctx context.Context, endpoint string) (body []byte, err error) {
	client.authToken, err = client.authenticate(ctx)

	if err == nil {

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-Plex-Token", client.authToken)

		var timer *prometheus.Timer
		if client.Options.PrometheusSummary != nil {
			timer = prometheus.NewTimer(client.Options.PrometheusSummary.WithLabelValues("plex", endpoint))
		}

		var resp *http.Response
		resp, err = client.Client.Do(req)

		if timer != nil {
			timer.ObserveDuration()
		}

		if err == nil {
			if resp.StatusCode == http.StatusOK {
				body, err = ioutil.ReadAll(resp.Body)
			} else {
				err = fmt.Errorf("%s", resp.Status)
			}
			_ = resp.Body.Close()
		}
	}

	return
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (client *Client) authenticate(ctx context.Context) (authToken string, err error) {
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

	var timer *prometheus.Timer
	if client.Options.PrometheusSummary != nil {
		timer = prometheus.NewTimer(client.Options.PrometheusSummary.WithLabelValues("plex", "auth"))
	}

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if timer != nil {
		timer.ObserveDuration()
	}

	if err == nil {
		if resp.StatusCode == http.StatusCreated {
			// TODO: there's three different places in the response where the authToken appears.
			// Which is the officially supported version?
			var authResponse struct {
				XMLName             xml.Name `xml:"user"`
				AuthenticationToken string   `xml:"authenticationToken,attr"`
			}

			body, _ := ioutil.ReadAll(resp.Body)
			if err = xml.Unmarshal(body, &authResponse); err == nil {
				authToken = authResponse.AuthenticationToken
			}
		} else {
			err = fmt.Errorf("%s", resp.Status)
		}
		_ = resp.Body.Close()
	}

	log.WithFields(log.Fields{
		"err":       err,
		"authToken": authToken,
	}).Debug("plex authenticate")

	return
}
