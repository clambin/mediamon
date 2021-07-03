package connectivity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/clambin/mediamon/internal/metrics"
)

// Probe to ping ipinfo.io
type Probe struct {
	Client *http.Client
	Token  string
}

// NewProbe creates a new Probe
func NewProbe(proxy *url.URL, token string) *Probe {
	return &Probe{
		Client: &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
			Timeout:   time.Second * 10,
		},
		Token: token,
	}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run(_ context.Context) (err error) {
	var connected float64

	if err = probe.getResponse(); err == nil {
		connected = 1
	} else {
		log.WithField("err", err).Warning("connectivity probe failed")
	}

	metrics.OpenVPNClientStatus.Set(connected)

	return
}

func (probe *Probe) getResponse() (err error) {
	var (
		response struct {
			IP       string
			Hostname string
			City     string
			Region   string
			Country  string
			Loc      string
			Org      string
			Postal   string
			Timezone string
		}
	)

	req, _ := http.NewRequest(http.MethodGet, "https://ipinfo.io/", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("token", probe.Token)
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	if resp, err = probe.Client.Do(req); err == nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		var body []byte
		if resp.StatusCode == 200 {
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				decoder := json.NewDecoder(bytes.NewReader(body))
				err = decoder.Decode(&response)
			}
		} else {
			err = errors.New(resp.Status)
		}
	}

	log.WithFields(log.Fields{"err": err, "response": response}).Debug("connectivity getResponse")

	return
}
