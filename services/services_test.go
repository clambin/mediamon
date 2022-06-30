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
	}{
		{filename: "testdata/complete.yaml", pass: true},
		{filename: "testdata/partial.yaml", pass: true},
		{filename: "testdata/invalid.yaml", pass: false},
		{filename: "testdata/invalid_proxy.yaml", pass: false},
	}

	for _, tt := range testCases {
		cfg, err := services.ParseConfigFile(tt.filename)
		if tt.pass == false {
			assert.Error(t, err, tt.filename)
			continue
		}
		require.NoError(t, err, tt.filename)
		assert.Equal(t, 5*time.Minute, cfg.OpenVPN.Connectivity.Interval)

		var body, golden []byte
		body, err = yaml.Marshal(cfg)

		gp := filepath.Join("testdata", t.Name()+"-"+slug.Make(tt.filename)+".golden")
		if *update {
			require.NoError(t, err, tt.filename)
			err = os.WriteFile(gp, body, 0644)
			require.NoError(t, err, tt.filename)

		}

		golden, err = os.ReadFile(gp)
		require.NoError(t, err, tt.filename)
		assert.Equal(t, string(golden), string(body), tt.filename)
	}
}

func TestParseMissingConfig(t *testing.T) {
	_, err := services.ParseConfigFile("not_a_file")
	require.Error(t, err)
}
