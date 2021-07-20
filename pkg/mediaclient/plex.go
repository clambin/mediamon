package mediaclient

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
	"strconv"
)

// PlexAPI interface
type PlexAPI interface {
	GetVersion(context.Context) (string, error)
	GetSessions(context.Context) (map[string]int, map[string]int, int, float64, error)
}

// PlexClient calls the Plex APIs
type PlexClient struct {
	Client    *http.Client
	URL       string
	Options   PlexOpts
	AuthURL   string
	UserName  string
	Password  string
	authToken string
}

// PlexOpts contains options to alter PlexClient behaviour
type PlexOpts struct {
	PrometheusSummary *prometheus.SummaryVec
}

// GetVersion retrieves the version of the Plex server
func (client *PlexClient) GetVersion(ctx context.Context) (version string, err error) {
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

type sessionStats struct {
	MediaContainer struct {
		Metadata []struct {
			User struct {
				Title string
			}
			TranscodeSession struct {
				Throttled     bool
				Speed         string
				VideoDecision string
			}
		}
	}
}

// GetSessions retrieves session information from the server.
//
// Returns:
//
//   - number of sessions per user
//   - number of sessions per type of decoder (direct/copy/transcode)
//   - number of active transcoders
//   - total transcoding speed
//   - any encounters errors
func (client *PlexClient) GetSessions(ctx context.Context) (users map[string]int, modes map[string]int, transcoding int, speed float64, err error) {
	var resp []byte
	if resp, err = client.call(ctx, "/status/sessions"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))

		var stats sessionStats
		err = decoder.Decode(&stats)

		if err == nil {
			users, modes, transcoding, speed, err = parseSessions(stats)
		}
	}

	log.WithFields(log.Fields{
		"err":         err,
		"users":       users,
		"modes":       modes,
		"transcoding": transcoding,
		"speed":       speed,
	}).Debug("plex getSessions")

	return
}

func parseSessions(stats sessionStats) (users map[string]int, modes map[string]int, transcoding int, speed float64, err error) {
	users = make(map[string]int)
	modes = make(map[string]int)

	for _, entry := range stats.MediaContainer.Metadata {
		// User sessions
		count, _ := users[entry.User.Title]
		users[entry.User.Title] = count + 1

		// Transcoders
		videoDecision := entry.TranscodeSession.VideoDecision
		if videoDecision == "" {
			videoDecision = "direct"
		}
		count, _ = modes[videoDecision]
		modes[videoDecision] = count + 1

		// Active transcoders
		if !entry.TranscodeSession.Throttled {
			transcoding++
			if entry.TranscodeSession.Speed != "" {
				var parsedSpeed float64
				if parsedSpeed, err = strconv.ParseFloat(entry.TranscodeSession.Speed, 64); err == nil {
					speed += parsedSpeed
				} else {
					log.WithFields(log.Fields{
						"err":   err,
						"Speed": entry.TranscodeSession.Speed,
					}).Warning("cannot parse Plex encoding speed")
				}
			}
		}
	}

	return
}

// call the specified Plex API endpoint
func (client *PlexClient) call(ctx context.Context, endpoint string) (body []byte, err error) {
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
func (client *PlexClient) authenticate(ctx context.Context) (authToken string, err error) {
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
