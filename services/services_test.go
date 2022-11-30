package services_test

import (
	"flag"
	"github.com/clambin/mediamon/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
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
		{filename: "complete", pass: true},
		{filename: "partial", pass: true},
		{filename: "invalid", pass: false},
		{filename: "envvar", pass: true, env: EnvVars{"PLEX_PASSWORD": "some-password"}},
		{filename: "invalid_proxy", pass: false},
	}

	for _, tt := range testCases {
		t.Run(tt.filename, func(t *testing.T) {
			require.NoError(t, tt.env.Set())

			cfg, err := services.ParseConfigFile(filepath.Join("testdata", tt.filename+".yaml"))
			if tt.pass == false {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, 5*time.Minute, cfg.OpenVPN.Connectivity.Interval)

			var body, golden []byte
			body, err = yaml.Marshal(cfg)
			require.NoError(t, err)

			gp := filepath.Join("testdata", strings.ToLower(t.Name())+".golden")
			if *update {
				require.NoError(t, os.WriteFile(gp, body, 0644))
			}

			golden, err = os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(golden), string(body))

			require.NoError(t, tt.env.Clear())
		})
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
