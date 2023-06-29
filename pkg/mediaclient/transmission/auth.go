package transmission

import (
	"bytes"
	"io"
	"net/http"
)

var _ http.RoundTripper = &authenticator{}

type authenticator struct {
	sessionID string
	next      http.RoundTripper
}

const transmissionSessionIDHeader = "X-Transmission-Session-Id"

func (a *authenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	var bodyCopy bytes.Buffer
	request.Body = io.NopCloser(io.TeeReader(request.Body, &bodyCopy))

	request.Header.Set(transmissionSessionIDHeader, a.sessionID)
	resp, err := a.next.RoundTrip(request)

	if err != nil || resp.StatusCode != http.StatusConflict {
		return resp, err
	}

	a.sessionID = resp.Header.Get(transmissionSessionIDHeader)
	request.Header.Set(transmissionSessionIDHeader, a.sessionID)
	request.Body = io.NopCloser(&bodyCopy)

	return a.next.RoundTrip(request)
}
