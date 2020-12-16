package plex

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mediamon/internal/metrics"
	"net/http"
	"testing"
)

func TestPlexAuthentication(t *testing.T) {
	apiClient := NewAPIClient(makeClient(), "", "user@example.com", "somepassword")
	assert.True(t, apiClient.authenticate())

	apiClient = NewAPIClient(makeClient(), "", "user@example.com", "badpassword")
	assert.False(t, apiClient.authenticate())
}

func TestProbe_Run(t *testing.T) {
	probe := NewProbeWithHTTPClient(makeClient(), "", "user@example.com", "somepassword")

	log.SetLevel(log.DebugLevel)
	probe.Run()

	// Test version
	_, ok := metrics.LoadValue("version", "plex", "SomeVersion")
	assert.True(t, ok)

	_, ok = metrics.LoadValue("version", "plex", "NotSomeVersion")
	assert.False(t, ok)

	// Test user count
	for _, testCase := range []struct {
		user  string
		ok    bool
		value float64
	}{
		{"foo", true, float64(1)},
		{"bar", true, float64(1)},
		{"snafu", true, float64(1)},
		{"arya", false, float64(-1)},
	} {
		userCount, ok := metrics.LoadValue("plex_session_count", testCase.user)
		assert.Equal(t, testCase.ok, ok, testCase.user)
		if ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}

	// Test transcoder
	for _, testCase := range []struct {
		user  string
		ok    bool
		value float64
	}{
		{"direct", true, float64(1)},
		{"copy", true, float64(1)},
		{"transcode", true, float64(1)},
		{"snafu", false, float64(-1)},
	} {
		modeCount, ok := metrics.LoadValue("plex_transcoder_type_count", testCase.user)
		assert.Equal(t, testCase.ok, ok, testCase.user)
		if ok {
			assert.Equal(t, testCase.value, modeCount, testCase.user)
		}
	}

	// Active transcoders
	encodingCount, ok := metrics.LoadValue("plex_transcoder_encoding_count", "plex")
	assert.True(t, ok)
	assert.Equal(t, float64(2), encodingCount)

	// Total encoding speed
	encodingSpeed, ok := metrics.LoadValue("plex_transcoder_speed_total", "plex")
	assert.True(t, ok)
	assert.Equal(t, float64(3.1), encodingSpeed)
}

// Stubbing the API Call

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// Responses

const (
	authResponse = `<?xml version="1.0" encoding="UTF-8"?>
<user email="user@example.com" id="1" uuid="1" username="user" authenticationToken="some_token" authToken="some_token">
  <subscription active="0" status="Inactive" plan=""></subscription>
  <entitlements all="0"></entitlements>
  <profile_settings/>
  <providers></providers>
  <services></services>
  <username>user</username>
  <email>user@example.com</email>
  <joined-at type="datetime">2000-01-01 00:00:00 UTC</joined-at>
  <authentication-token>some_token</authentication-token>
</user>`

	identityResponse = `{
  "MediaContainer": {
    "size": 0,
    "claimed": true,
    "machineIdentifier": "SomeUUID",
    "version": "SomeVersion"
  }
}`

	sessionsResponse = `
{
  "MediaContainer": 
  {
    "size": 2,
    "Metadata": 
    [
      {
        "User": 
        {
          "title": "foo"
        }
      },
      {
        "User":
        {
          "title": "bar"
        },
        "TranscodeSession":
        {
          "throttled": false,
          "speed": "3.1",
          "videoDecision": "copy"
        }
      },
      {
        "User":
        {
          "title": "snafu"
        },
        "TranscodeSession":
        {
          "throttled": true,
          "speed": "3.1",
          "videoDecision": "transcode"
        }
      }
    ]
  }
}`
)

// makeClient returns a stubbed covid.APIClient
func makeClient() *http.Client {
	header := make(http.Header)

	return NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://plex.tv/users/sign_in.xml" {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return &http.Response{
					StatusCode: 500,
					Status:     err.Error(),
					Header:     header,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}
			}
			if string(body) == `user%5Blogin%5D=user@example.com&user%5Bpassword%5D=somepassword` {
				return &http.Response{
					StatusCode: 201,
					Header:     header,
					Body:       ioutil.NopCloser(bytes.NewBufferString(authResponse)),
				}
			}
			return &http.Response{
				StatusCode: 401,
				Status:     "Unauthorized",
				Header:     header,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}

		} else if req.URL.Path == "/identity" {
			return &http.Response{
				StatusCode: 200,
				Header:     header,
				Body:       ioutil.NopCloser(bytes.NewBufferString(identityResponse)),
			}
		} else if req.URL.Path == "/status/sessions" {
			return &http.Response{
				StatusCode: 200,
				Header:     header,
				Body:       ioutil.NopCloser(bytes.NewBufferString(sessionsResponse)),
			}
		}
		return &http.Response{
			StatusCode: 500,
			Status:     "Not implemented",
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	})
}
