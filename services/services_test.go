package services_test

import (
	"flag"
	"github.com/clambin/mediamon/services"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update .golden files")

func TestParseConfigFile(t *testing.T) {
	testCases := []struct {
		filename string
		pass     bool
		env      EnvVars
	}{
		{filename: "testdata/complete.yaml", pass: true},
		{filename: "testdata/partial.yaml", pass: true},
		{filename: "testdata/invalid.yaml", pass: false},
		{filename: "testdata/envvar.yaml", pass: true, env: EnvVars{"PLEX_PASSWORD": "some-password"}},
		{filename: "testdata/invalid_proxy.yaml", pass: false},
	}

	for _, tt := range testCases {
		err := tt.env.Set()
		require.NoError(t, err)

		var cfg *services.Config
		cfg, err = services.ParseConfigFile(tt.filename)
		if tt.pass == false {
			assert.Error(t, err, tt.filename)
			continue
		}
		require.NoError(t, err, tt.filename)
		assert.Equal(t, 5*time.Minute, cfg.OpenVPN.Connectivity.Interval)

		var body, golden []byte
		body, err = yaml.Marshal(cfg)
		require.NoError(t, err, tt.filename)

		gp := filepath.Join("testdata", t.Name()+"-"+slug.Make(tt.filename)+".golden")
		if *update {
			err = os.WriteFile(gp, body, 0644)
			require.NoError(t, err, tt.filename)
		}

		golden, err = os.ReadFile(gp)
		require.NoError(t, err, tt.filename)
		assert.Equal(t, string(golden), string(body), tt.filename)

		err = tt.env.Clear()
		require.NoError(t, err)
	}
}

func TestParseMissingConfig(t *testing.T) {
	_, err := services.ParseConfigFile("not_a_file")
	require.Error(t, err)
}

type EnvVars map[string]string

func (e EnvVars) Set() error {
	for key, value := range e {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (e EnvVars) Clear() error {
	for key := range e {
		if err := os.Unsetenv(key); err != nil {
			return err
		}
	}
	return nil
}
