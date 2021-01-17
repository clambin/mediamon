package plex

import (
	log "github.com/sirupsen/logrus"
	"mediamon/internal/metrics"
	"mediamon/pkg/mediaclient"
	"net/http"
)

// Probe measures Plex metrics
type Probe struct {
	mediaclient.PlexAPI
	Users map[string]int
	Modes map[string]int
}

// NewProbe creates a new Probe
func NewProbe(url, username, password string) *Probe {
	return &Probe{
		&mediaclient.PlexClient{
			Client:   &http.Client{},
			URL:      url,
			UserName: username,
			Password: password,
		},
		make(map[string]int),
		make(map[string]int),
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run() error {
	var (
		err         error
		version     string
		users       map[string]int
		modes       map[string]int
		transcoding int
		speed       float64
	)
	// Get the version
	if version, err = probe.GetVersion(); err == nil {
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
	if users, modes, transcoding, speed, err = probe.GetSessions(); err == nil {
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

	return err
}
