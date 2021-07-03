package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/gotools/metrics"
	"github.com/clambin/mediamon/internal/plex"
	"github.com/stretchr/testify/assert"
	"testing"
)

type client struct {
	failing  bool
	scenario int
}

func (client *client) GetVersion(_ context.Context) (string, error) {
	if client.failing {
		return "", errors.New("failing")
	}
	return "foo", nil
}

func (client *client) GetSessions(_ context.Context) (map[string]int, map[string]int, int, float64, error) {
	if client.failing {
		return nil, nil, 0, 0.0, errors.New("failing")
	}
	users := make(map[string]int)
	modes := make(map[string]int)

	users["foo"] = 1
	users["bar"] = 1
	modes["direct"] = 1
	modes["copy"] = 1
	if client.scenario == 0 {
		users["snafu"] = 2
		modes["transcode"] = 2
	}
	return users, modes, 2, 3.1, nil
}

func TestProbe_Run(t *testing.T) {
	probe := plex.Probe{
		PlexAPI: &client{},
		Users:   make(map[string]int),
		Modes:   make(map[string]int),
	}

	// log.SetLevel(log.DebugLevel)
	_ = probe.Run(context.Background())

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
	c := client{}
	probe := plex.Probe{
		PlexAPI: &c,
		Users:   make(map[string]int),
		Modes:   make(map[string]int),
	}

	// log.SetLevel(log.DebugLevel)
	_ = probe.Run(context.Background())

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
	c.scenario = 1

	// Run again
	_ = probe.Run(context.Background())

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

func TestProbe_Fail(t *testing.T) {
	c := client{failing: true}
	probe := plex.Probe{
		PlexAPI: &c,
		Users:   make(map[string]int),
		Modes:   make(map[string]int),
	}

	err := probe.Run(context.Background())
	assert.NotNil(t, err)
}
