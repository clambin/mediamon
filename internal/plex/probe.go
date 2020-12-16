package plex

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"

	"mediamon/internal/metrics"
	"net/http"
)

type Probe struct {
	apiClient *APIClient
}

func NewProbe(url, username, password string) *Probe {
	return NewProbeWithHTTPClient(&http.Client{}, url, username, password)
}

func NewProbeWithHTTPClient(client *http.Client, url, username, password string) *Probe {
	return &Probe{apiClient: NewAPIClient(client, url, username, password)}
}

func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err != nil {
		log.Warningf("Could not get Plex version: %s", err)
	} else {
		metrics.Publish("version", 1, "plex", version)
	}

	// Get sessions
	// FIXME: need to maintain a list of all users reported and report 0 if they're not in the current measurement
	users, modes, transcoding, speed, err := probe.getSessions()

	if err == nil {
		for user, value := range users {
			metrics.Publish("plex_session_count", float64(value), user)
		}
		for mode, value := range modes {
			metrics.Publish("plex_transcoder_type_count", float64(value), mode)
		}
		metrics.Publish("plex_transcoder_encoding_count", float64(transcoding), "plex")
		metrics.Publish("plex_transcoder_speed_total", speed, "plex")
	}

}

func (probe *Probe) getVersion() (string, error) {
	var stats = struct {
		MediaContainer struct {
			Version string
		}
	}{}

	resp, err := probe.apiClient.Call("/identity")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			log.Debugf("API call result: %v", stats)
			return stats.MediaContainer.Version, nil
		}
	}
	return "", err
}

func (probe *Probe) getSessions() (map[string]int, map[string]int, int, float64, error) {
	var stats = struct {
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
	}{}

	resp, err := probe.apiClient.Call("/status/sessions")
	if err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
		if err == nil {
			log.Debugf("API call result: %v", stats)

			users := make(map[string]int, 0)
			modes := make(map[string]int, 0)
			transcoding := 0
			speed := float64(0)

			for _, entry := range stats.MediaContainer.Metadata {
				// User sessions
				userCount, ok := users[entry.User.Title]
				if ok == false {
					users[entry.User.Title] = 1
				} else {
					users[entry.User.Title] = userCount + 1
				}
				// Transcoders
				videoDecision := entry.TranscodeSession.VideoDecision
				if videoDecision == "" {
					videoDecision = "direct"
				}
				modeCount, ok := modes[videoDecision]
				if ok == false {
					modes[videoDecision] = 1
				} else {
					modes[videoDecision] = modeCount + 1
				}

				// Active transcoders

				if !entry.TranscodeSession.Throttled {
					transcoding++
					if entry.TranscodeSession.Speed != "" {
						if parsedSpeed, err := strconv.ParseFloat(entry.TranscodeSession.Speed, 64); err != nil {
							log.Warningf("cannot parse Plex encoding speed: '%s'", entry.TranscodeSession.Speed)
						} else {
							speed += parsedSpeed
						}
					}
				}
			}

			return users, modes, transcoding, speed, nil
		}
	}
	return nil, nil, -1, float64(-1), err
}
