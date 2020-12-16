package plex

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"

	"mediamon/internal/metrics"
	"net/http"
)

// Probe to measure Plex metrics
type Probe struct {
	apiClient *Client
	users     map[string]int
	modes     map[string]int
}

// NewProbe creates a new Probe
func NewProbe(url, username, password string) *Probe {
	return NewProbeWithHTTPClient(&http.Client{}, url, username, password)
}

// NewProbeWithHTTPClient creates a probe with a specified http.Client
// Used to stub API calls during unit testing
func NewProbeWithHTTPClient(client *http.Client, url, username, password string) *Probe {
	return &Probe{
		apiClient: NewAPIWithHTTPClient(client, url, username, password),
		users:     make(map[string]int),
		modes:     make(map[string]int),
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err != nil {
		log.Warningf("Could not get Plex version: %s", err)
	} else {
		metrics.Publish("version", 1, "plex", version)
	}

	// Reset current statistics
	for user := range probe.users {
		probe.users[user] = 0
	}
	for mode := range probe.modes {
		probe.modes[mode] = 0
	}

	// Get sessions
	users, modes, transcoding, speed, err := probe.getSessions()

	if err == nil {
		// Update statistics
		for user, count := range users {
			if oldCount, ok := probe.users[user]; ok {
				probe.users[user] = oldCount + count
			} else {
				probe.users[user] = count
			}
		}
		for mode, count := range modes {
			if oldCount, ok := probe.modes[mode]; ok {
				probe.modes[mode] = oldCount + count
			} else {
				probe.modes[mode] = count
			}
		}

		// Report
		for user, value := range probe.users {
			metrics.Publish("plex_session_count", float64(value), user)
		}
		for mode, value := range probe.modes {
			metrics.Publish("plex_transcoder_type_count", float64(value), mode)
		}
		metrics.Publish("plex_transcoder_encoding_count", float64(transcoding))
		metrics.Publish("plex_transcoder_speed_total", speed)
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
