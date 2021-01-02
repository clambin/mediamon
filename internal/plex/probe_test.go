package plex_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/clambin/httpstub"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/plex"
	"mediamon/pkg/metrics"
)

func TestProbe_Run(t *testing.T) {
	probe := plex.NewProbeWithHTTPClient(httpstub.NewTestClient(loopback), "", "user@example.com", "somepassword")

	// log.SetLevel(log.DebugLevel)
	probe.Run()

	// Test version
	_, err := metrics.LoadValue("mediaserver_server_info", "plex", "SomeVersion")
	assert.Nil(t, err)

	_, err = metrics.LoadValue("mediaserver_server_info", "plex", "NotSomeVersion")
	assert.Nil(t, err)

	// Test user count
	for _, testCase := range []struct {
		user  string
		ok    bool
		value float64
	}{
		{"foo", true, 1.0},
		{"bar", true, 1.0},
		{"snafu", true, 2.0},
		// Current implementation of LoadValue can't detect missing labels. Will return 0
		{"ufans", true, 0.0},
	} {
		userCount, err := metrics.LoadValue("mediaserver_plex_session_count", testCase.user)
		if testCase.ok {
			assert.Nil(t, err, testCase.user)
			assert.Equal(t, testCase.value, userCount, testCase.user)
		} else {
			assert.NotNil(t, err, testCase.user)
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
		// Current implementation of LoadValue can't detect missing labels. Will return 0
		{"snafu", true, 0.0},
	} {
		modeCount, err := metrics.LoadValue("mediaserver_plex_transcoder_type_count", testCase.mode)
		if testCase.ok {
			assert.Nil(t, err, testCase.mode)
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		} else {
			assert.NotNil(t, err, testCase.mode)
		}
	}

	// Active transcoders
	encodingCount, err := metrics.LoadValue("mediaserver_plex_transcoder_encoding_count")
	assert.Nil(t, err)
	assert.Equal(t, 2.0, encodingCount)

	// Total encoding speed
	encodingSpeed, err := metrics.LoadValue("mediaserver_plex_transcoder_speed_total")
	assert.Nil(t, err)
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
		userCount, err := metrics.LoadValue("mediaserver_plex_session_count", testCase.user)
		if testCase.ok {
			assert.Nil(t, err, testCase.user)
			assert.Equal(t, testCase.value, userCount, testCase.user)
		} else {
			assert.NotNil(t, err, testCase.user)
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
		modeCount, err := metrics.LoadValue("mediaserver_plex_transcoder_type_count", testCase.mode)
		if testCase.ok {
			assert.Nil(t, err, testCase.mode)
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		} else {
			assert.NotNil(t, err, testCase.mode)
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
		userCount, err := metrics.LoadValue("mediaserver_plex_session_count", testCase.user)
		if testCase.ok {
			assert.Nil(t, err, testCase.user)
			assert.Equal(t, testCase.value, userCount, testCase.user)
		} else {
			assert.NotNil(t, err, testCase.user)
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
		modeCount, err := metrics.LoadValue("mediaserver_plex_transcoder_type_count", testCase.mode)
		if testCase.ok {
			assert.Nil(t, err, testCase.mode)
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		} else {
			assert.NotNil(t, err, testCase.mode)
		}
	}
}

func TestFailingServer(t *testing.T) {
	probe := plex.NewProbeWithHTTPClient(
		httpstub.NewTestClient(httpstub.Failing),
		"http://example.com",
		"user@example.com",
		"password",
	)

	assert.NotPanics(t, func() { probe.Run() })
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
