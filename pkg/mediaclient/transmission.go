package mediaclient

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// TransmissionAPI interface
type TransmissionAPI interface {
	GetVersion() (string, error)
	GetStats() (int, int, int, int, error)
}

// TransmissionClient calls the Transmission APIs
type TransmissionClient struct {
	Client    *http.Client
	URL       string
	SessionID string
}

// GetVersion determines the version of the Transmission server
func (client *TransmissionClient) GetVersion() (version string, err error) {
	var (
		resp  []byte
		stats = struct {
			Arguments struct {
				Version string
			}
		}{}
	)

	if resp, err = client.call("session-get"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		if err = decoder.Decode(&stats); err == nil {
			version = stats.Arguments.Version
		}
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": stats.Arguments.Version,
	}).Debug("transmission GetVersion")

	return
}

// GetStats gets torrent & up/download stats from Transmission.
//
// Returns:
//   - active torrents
//   - paused torrents
//   - total download speed
//   - total upload speed
//   - encountered error
func (client *TransmissionClient) GetStats() (active int, paused int, download int, upload int, err error) {
	var (
		resp  []byte
		stats = struct {
			Arguments struct {
				ActiveTorrentCount int
				PausedTorrentCount int
				UploadSpeed        int
				DownloadSpeed      int
			}
			Result string
		}{}
	)

	if resp, err = client.call("session-stats"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		if err = decoder.Decode(&stats); err == nil {
			active = stats.Arguments.ActiveTorrentCount
			paused = stats.Arguments.PausedTorrentCount
			download = stats.Arguments.DownloadSpeed
			upload = stats.Arguments.UploadSpeed
		}
	}

	log.WithFields(log.Fields{
		"err":   err,
		"stats": stats.Arguments,
	}).Debug("transmission GetStats")

	return
}

// call the specified Transmission API endpoint
func (client *TransmissionClient) call(method string) (response []byte, err error) {
	for err == nil {
		if client.SessionID, err = client.getSessionID(); err == nil {
			req, _ := http.NewRequest(http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Transmission-Session-Id", client.SessionID)

			var resp *http.Response
			if resp, err = client.Client.Do(req); err == nil {

				if resp.StatusCode == 409 {
					// Transmission-Session-Id has expired. Get the new one and retry
					client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
					resp.Body.Close()
				} else {
					if resp.StatusCode == 200 {
						response, err = ioutil.ReadAll(resp.Body)
						break
					} else {
						err = errors.New(resp.Status)
					}
					resp.Body.Close()
				}
			}
		}
	}
	return
}

func (client *TransmissionClient) getSessionID() (sessionID string, err error) {
	sessionID = client.SessionID

	if sessionID == "" {
		req, _ := http.NewRequest(http.MethodPost, client.URL, bytes.NewBufferString("{ \"method\": \"session-get\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

		var resp *http.Response
		if resp, err = client.Client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 409 || resp.StatusCode == 200 {
				sessionID = resp.Header.Get("X-Transmission-Session-Id")
			}
		}
		log.WithFields(log.Fields{
			"err":       err,
			"sessionID": sessionID,
		}).Debug("transmission getSessionID")
	}

	return
}
