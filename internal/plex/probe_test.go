package plex_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"mediamon/internal/metrics"
	"mediamon/internal/plex"

	"mediamon/internal/httpstub"
)

func TestProbe_Run(t *testing.T) {
	probe := plex.NewProbeWithHTTPClient(httpstub.NewTestClient(loopback), "", "user@example.com", "somepassword")

	// log.SetLevel(log.DebugLevel)
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
		{"foo", true, 1.0},
		{"bar", true, 1.0},
		{"snafu", true, 2.0},
		{"ufans", false, -1.0},
	} {
		userCount, ok := metrics.LoadValue("plex_session_count", testCase.user)
		assert.Equal(t, testCase.ok, ok, testCase.user)
		if ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}

	// Test transcoder
	for _, testCase := range []struct {
		mode  string
		ok    bool
		value float64
	}{
		{"direct", true, 1.0},
		{"copy", true, 1.0},
		{"transcode", true, 2.0},
		{"snafu", false, -1.0},
	} {
		modeCount, ok := metrics.LoadValue("plex_transcoder_type_count", testCase.mode)
		assert.Equal(t, testCase.ok, ok, testCase.mode)
		if ok {
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		}
	}

	// Active transcoders
	encodingCount, ok := metrics.LoadValue("plex_transcoder_encoding_count")
	assert.True(t, ok)
	assert.Equal(t, 2.0, encodingCount)

	// Total encoding speed
	encodingSpeed, ok := metrics.LoadValue("plex_transcoder_speed_total")
	assert.True(t, ok)
	assert.Equal(t, 3.1, encodingSpeed)
}

func TestCachedUsers(t *testing.T) {
	probe := plex.NewProbeWithHTTPClient(httpstub.NewTestClient(loopback), "", "user@example.com", "somepassword")

	// log.SetLevel(log.DebugLevel)
	probe.Run()

	// Test user count
	for _, testCase := range []struct {
		user  string
		ok    bool
		value float64
	}{
		{"foo", true, 1.0},
		{"bar", true, 1.0},
		{"snafu", true, 2.0},
	} {
		userCount, ok := metrics.LoadValue("plex_session_count", testCase.user)
		assert.Equal(t, testCase.ok, ok, testCase.user)
		if ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}

	// Test transcoder
	for _, testCase := range []struct {
		mode  string
		ok    bool
		value float64
	}{
		{"direct", true, 1.0},
		{"copy", true, 1.0},
		{"transcode", true, 2.0},
	} {
		modeCount, ok := metrics.LoadValue("plex_transcoder_type_count", testCase.mode)
		assert.Equal(t, testCase.ok, ok, testCase.mode)
		if ok {
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		}
	}

	// Switch to response w/out Snafu user & related transcoders
	sessionsResponse = sessionsResponse2

	// Run again
	probe.Run()

	// Snafu should still be reported but with zero sessions
	for _, testCase := range []struct {
		user  string
		ok    bool
		value float64
	}{
		{"foo", true, 1.0},
		{"bar", true, 1.0},
		{"snafu", true, 0.0},
	} {
		userCount, ok := metrics.LoadValue("plex_session_count", testCase.user)
		assert.Equal(t, testCase.ok, ok, testCase.user)
		if ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}

	// Same for transcode mode
	for _, testCase := range []struct {
		mode  string
		ok    bool
		value float64
	}{
		{"direct", true, 1.0},
		{"copy", true, 1.0},
		{"transcode", true, 0.0},
	} {
		modeCount, ok := metrics.LoadValue("plex_transcoder_type_count", testCase.mode)
		assert.Equal(t, testCase.ok, ok, testCase.mode)
		if ok {
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		}
	}
}

// Server loopback function

func loopback(req *http.Request) *http.Response {
	if req.URL.String() == "https://plex.tv/users/sign_in.xml" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return &http.Response{
				StatusCode: 500,
				Status:     err.Error(),
				Header:     nil,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
		}
		if string(body) == `user%5Blogin%5D=user@example.com&user%5Bpassword%5D=somepassword` {
			return &http.Response{
				StatusCode: 201,
				Header:     nil,
				Body:       ioutil.NopCloser(bytes.NewBufferString(authResponse)),
			}
		}
		return &http.Response{
			StatusCode: 401,
			Status:     "Unauthorized",
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}

	} else if req.URL.Path == "/identity" {
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(identityResponse)),
		}
	} else if req.URL.Path == "/status/sessions" {
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(sessionsResponse)),
		}
	}
	return &http.Response{
		StatusCode: 500,
		Status:     "Not implemented",
		Header:     nil,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}
}

// Responses

var (
	// default session response answer
	sessionsResponse = sessionsResponse1
)

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

	sessionsResponse1 = `
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

	sessionsResponse2 = `
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
      }
    ]
  }
}`
)
