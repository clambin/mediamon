package mediaclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/clambin/mediamon/internal/version"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

// PlexAPI interface
type PlexAPI interface {
	GetVersion() (string, error)
	GetSessions() (map[string]int, map[string]int, int, float64, error)
}

// PlexClient calls the Plex APIs
type PlexClient struct {
	Client    *http.Client
	URL       string
	UserName  string
	Password  string
	authToken string
}

// GetVersion retrieves the version of the Plex server
func (client *PlexClient) GetVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats struct {
			MediaContainer struct {
				Version string
			}
		}
	)

	if resp, err = client.call("/identity"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": stats.MediaContainer.Version,
	}).Debug("plex getVersion")

	return stats.MediaContainer.Version, err
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
func (client *PlexClient) GetSessions() (map[string]int, map[string]int, int, float64, error) {
	var (
		err         error
		resp        []byte
		users       = make(map[string]int, 0)
		modes       = make(map[string]int, 0)
		transcoding int
		speed       float64
		ok          bool
		count       int
		stats       struct {
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
	)

	if resp, err = client.call("/status/sessions"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		if err = decoder.Decode(&stats); err == nil {
			for _, entry := range stats.MediaContainer.Metadata {
				// User sessions
				count, ok = users[entry.User.Title]
				if ok {
					users[entry.User.Title] = count + 1
				} else {
					users[entry.User.Title] = 1
				}
				// Transcoders
				videoDecision := entry.TranscodeSession.VideoDecision
				if videoDecision == "" {
					videoDecision = "direct"
				}
				count, ok = modes[videoDecision]
				if ok {
					modes[videoDecision] = count + 1
				} else {
					modes[videoDecision] = 1
				}
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
		}
	}

	log.WithFields(log.Fields{
		"err":         err,
		"users":       users,
		"modes":       modes,
		"transcoding": transcoding,
		"speed":       speed,
	}).Debug("plex getSessions")

	return users, modes, transcoding, speed, err
}

// call the specified Plex API endpoint
func (client *PlexClient) call(endpoint string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	if client.authToken, err = client.authenticate(); err == nil && client.authToken != "" {

		req, _ := http.NewRequest("GET", client.URL+endpoint, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("X-Plex-Token", client.authToken)

		if resp, err = client.Client.Do(req); err == nil {
			if resp.StatusCode == 200 {
				body, err = ioutil.ReadAll(resp.Body)
			} else {
				err = errors.New(resp.Status)
			}
			_ = resp.Body.Close()
		}
	}

	return body, err
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (client *PlexClient) authenticate() (string, error) {
	var (
		err       error
		resp      *http.Response
		authToken = client.authToken
	)

	if authToken == "" {
		authBody := fmt.Sprintf("user%%5Blogin%%5D=%s&user%%5Bpassword%%5D=%s",
			client.UserName,
			client.Password,
		)

		req, _ := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", bytes.NewBufferString(authBody))
		req.Header.Add("X-Plex-Product", "github.com/clambin/mediamon")
		req.Header.Add("X-Plex-Version", version.BuildVersion)
		req.Header.Add("X-Plex-Client-Identifier", "github.com/clambin/mediamon-v"+version.BuildVersion)

		if resp, err = client.Client.Do(req); err == nil {
			if resp.StatusCode == 201 {
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
				err = errors.New(resp.Status)
			}
			_ = resp.Body.Close()
		}
		log.WithFields(log.Fields{
			"err":       err,
			"authToken": authToken,
		}).Debug("plex authenticate")
	}

	return authToken, err
}
