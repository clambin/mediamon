package transmission

import (
	"bytes"
	"io"
	"net/http"
)

var _ http.RoundTripper = &authenticator{}

type authenticator struct {
	sessionId string
	next      http.RoundTripper
}

const transmissionSessionIDHeader = "X-Transmission-Session-Id"

func (a *authenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	var bodyCopy bytes.Buffer
	request.Body = io.NopCloser(io.TeeReader(request.Body, &bodyCopy))

	request.Header.Set(transmissionSessionIDHeader, a.sessionId)
	resp, err = a.next.RoundTrip(request)

	if err == nil && resp.StatusCode == http.StatusConflict {
		a.sessionId = resp.Header.Get(transmissionSessionIDHeader)
		request.Header.Set(transmissionSessionIDHeader, a.sessionId)
		request.Body = io.NopCloser(bytes.NewBuffer(bodyCopy.Bytes()))
		resp, err = a.next.RoundTrip(request)
	}
	return resp, err
}
