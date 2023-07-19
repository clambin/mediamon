package transmission

import (
	"bytes"
	"github.com/clambin/go-common/httpclient"
	"io"
	"net/http"
)

var _ http.RoundTripper = &authenticator{}

type authenticator struct {
	sessionID string
	next      http.RoundTripper
}

func withAuthenticator() httpclient.Option {
	return func(next http.RoundTripper) http.RoundTripper {
		a := authenticator{next: next}
		return &a
	}
}

const transmissionSessionIDHeader = "X-Transmission-Session-Id"

func (a *authenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	var bodyCopy bytes.Buffer
	origBody := io.TeeReader(request.Body, &bodyCopy)
	request.Body = io.NopCloser(origBody)

	request.Header.Set(transmissionSessionIDHeader, a.sessionID)
	resp, err := a.next.RoundTrip(request)

	if err != nil || resp.StatusCode != http.StatusConflict {
		return resp, err
	}

	_ = resp.Body.Close()
	a.sessionID = resp.Header.Get(transmissionSessionIDHeader)
	request.Header.Set(transmissionSessionIDHeader, a.sessionID)
	request.Body = io.NopCloser(&bodyCopy)

	return a.next.RoundTrip(request)
}
