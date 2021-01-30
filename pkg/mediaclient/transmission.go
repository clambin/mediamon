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
func (client *TransmissionClient) GetVersion() (string, error) {
	var (
		err   error
		resp  []byte
		stats = struct {
			Arguments struct {
				Version string
			}
		}{}
	)

	if resp, err = client.call("session-get"); err == nil {
		decoder := json.NewDecoder(bytes.NewReader(resp))
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":     err,
		"version": stats.Arguments.Version,
	}).Debug("transmission GetVersion")

	return stats.Arguments.Version, err
}

// GetStats gets torrent & up/download stats from Transmission.
//
// Returns:
//   - active torrents
//   - paused torrents
//   - total download speed
//   - total upload speed
//   - encountered error
func (client *TransmissionClient) GetStats() (int, int, int, int, error) {
	var (
		err   error
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
		err = decoder.Decode(&stats)
	}

	log.WithFields(log.Fields{
		"err":   err,
		"stats": stats.Arguments,
	}).Debug("transmission GetStats")

	return stats.Arguments.ActiveTorrentCount,
		stats.Arguments.PausedTorrentCount,
		stats.Arguments.DownloadSpeed,
		stats.Arguments.UploadSpeed,
		err
}

// call the specified Transmission API endpoint
func (client *TransmissionClient) call(method string) ([]byte, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)

	for {
		if client.SessionID, err = client.getSessionID(); err == nil {
			req, _ := http.NewRequest("POST", client.URL, bytes.NewBufferString("{ \"method\": \""+method+"\" }"))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Transmission-Session-Id", client.SessionID)

			if resp, err = client.Client.Do(req); err == nil {

				if resp.StatusCode == 409 {
					// Transmission-Session-Id has expired. Get the new one and retry
					client.SessionID = resp.Header.Get("X-Transmission-Session-Id")
					resp.Body.Close()
				} else {
					if resp.StatusCode == 200 {
						body, err = ioutil.ReadAll(resp.Body)
					} else {
						err = errors.New(resp.Status)
					}
					resp.Body.Close()
					break
				}
			}
		}
	}

	return body, err
}

func (client *TransmissionClient) getSessionID() (string, error) {
	var (
		err       error
		resp      *http.Response
		sessionID = client.SessionID
	)

	if sessionID == "" {
		req, _ := http.NewRequest("POST", client.URL, bytes.NewBufferString("{ \"method\": \"session-get\" }"))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Transmission-Session-Id", client.SessionID)

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

	return sessionID, err
}
