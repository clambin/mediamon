package plex

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/metrics"
)

// Probe to measure Plex metrics
type Probe struct {
	Client
	Users map[string]int
	Modes map[string]int
}

// NewProbe creates a new Probe
func NewProbe(url, username, password string) *Probe {
	return &Probe{
		Client{Client: &http.Client{}, URL: url, UserName: username, Password: password},
		make(map[string]int),
		make(map[string]int),
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() {
	// Get the version
	if version, err := probe.getVersion(); err == nil {
		metrics.MediaServerVersion.WithLabelValues("plex", version).Set(1)
	} else {
		log.WithField("err", err).Warning("Could not get Plex version")
	}

	// Reset current statistics
	for user := range probe.Users {
		probe.Users[user] = 0
	}
	for mode := range probe.Modes {
		probe.Modes[mode] = 0
	}

	// Get sessions
	if users, modes, transcoding, speed, err := probe.getSessions(); err == nil {
		// Update statistics
		for user, count := range users {
			if oldCount, ok := probe.Users[user]; ok {
				probe.Users[user] = oldCount + count
			} else {
				probe.Users[user] = count
			}
		}
		for mode, count := range modes {
			if oldCount, ok := probe.Modes[mode]; ok {
				probe.Modes[mode] = oldCount + count
			} else {
				probe.Modes[mode] = count
			}
		}

		// Report
		for user, value := range probe.Users {
			metrics.PlexSessionCount.WithLabelValues(user).Set(float64(value))
		}
		for mode, value := range probe.Modes {
			metrics.PlexTranscoderTypeCount.WithLabelValues(mode).Set(float64(value))
		}
		metrics.PlexTranscoderEncodingCount.Set(float64(transcoding))
		metrics.PlexTranscoderSpeedTotal.Set(speed)
	} else {
		log.WithField("err", err).Warning("could not get Plex sessions")
	}
}

func (probe *Probe) getVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats struct {
			MediaContainer struct {
				Version string
			}
		}
	)

	if resp, err = probe.call("/identity"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": stats.MediaContainer.Version,
	}).Debug("plex getVersion")

	return stats.MediaContainer.Version, err
}

func (probe *Probe) getSessions() (map[string]int, map[string]int, int, float64, error) {
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

	if resp, err = probe.call("/status/sessions"); err == nil {
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
