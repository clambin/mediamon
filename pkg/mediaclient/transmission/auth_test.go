package transmission

import (
	"github.com/clambin/go-common/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAuthenticator_RoundTrip(t *testing.T) {
	h := server{sessionID: "1234"}

	c := http.Client{Transport: httpclient.NewRoundTripper(
		withAuthenticator(),
		httpclient.WithRoundTripper(&h),
	)}

	req, _ := http.NewRequest(http.MethodPost, "application/json", io.NopCloser(strings.NewReader("foo")))
	resp, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "1234", resp.Header.Get(transmissionSessionIDHeader))

	// simulate an expired session ID
	h.sessionID = "4321"
	req, _ = http.NewRequest(http.MethodPost, "application/json", io.NopCloser(strings.NewReader("foo")))
	resp, err = c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "4321", resp.Header.Get(transmissionSessionIDHeader))
}

type server struct {
	sessionID string
}

func (s server) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := http.Response{StatusCode: http.StatusOK}
	resp.Header = make(http.Header)
	resp.Header.Set(transmissionSessionIDHeader, s.sessionID)
	body, _ := io.ReadAll(req.Body)

	if req.Header.Get(transmissionSessionIDHeader) != s.sessionID {
		resp.StatusCode = http.StatusConflict
		return &resp, nil
	}
	if string(body) != "foo" {
		resp.StatusCode = http.StatusBadRequest
	}
	return &resp, nil
}
