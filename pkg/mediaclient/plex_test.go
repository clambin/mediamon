package mediaclient_test

import (
	"bytes"
	"github.com/clambin/gotools/httpstub"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mediamon/pkg/mediaclient"
	"net/http"
	"testing"
)

func TestPlexClient_GetVersion(t *testing.T) {
	client := &mediaclient.PlexClient{
		Client:   httpstub.NewTestClient(plexLoopback),
		URL:      "",
		UserName: "user@example.com",
		Password: "somepassword",
	}

	version, err := client.GetVersion()
	assert.Nil(t, err)
	assert.Equal(t, "SomeVersion", version)
}

func TestPlexClient_GetStats(t *testing.T) {
	client := &mediaclient.PlexClient{
		Client:   httpstub.NewTestClient(plexLoopback),
		URL:      "",
		UserName: "user@example.com",
		Password: "somepassword",
	}

	users, modes, transcoding, speed, err := client.GetSessions()
	assert.Nil(t, err)
	assert.Len(t, users, 3)
	assert.Len(t, modes, 3)
	assert.Equal(t, 2, transcoding)
	assert.Equal(t, 3.1, speed)

	// User count
	for _, testCase := range []struct {
		user  string
		ok    bool
		value int
	}{
		{"foo", true, 1},
		{"bar", true, 1},
		{"snafu", true, 2},
		{"ufans", false, 0},
	} {
		userCount, ok := users[testCase.user]
		assert.Equal(t, testCase.ok, ok)
		if testCase.ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}
	// Mode count
	for _, testCase := range []struct {
		mode  string
		ok    bool
		value int
	}{
		{"direct", true, 1},
		{"copy", true, 1},
		{"transcode", true, 2},
		{"snafu", false, 0},
	} {
		modeCount, ok := modes[testCase.mode]
		assert.Equal(t, testCase.ok, ok)
		if testCase.ok {
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		}
	}
}

func TestPlexClient_Authentication(t *testing.T) {
	client := &mediaclient.PlexClient{
		Client:   httpstub.NewTestClient(plexLoopback),
		URL:      "",
		UserName: "user@example.com",
		Password: "badpassword",
	}

	_, err := client.GetVersion()
	assert.NotNil(t, err)
	assert.Equal(t, "Unauthorized", err.Error())
}

// Server loopback function
func plexLoopback(req *http.Request) *http.Response {
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
