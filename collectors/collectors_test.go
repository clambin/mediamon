package collectors_test

import (
	"bytes"
	"flag"
	"github.com/clambin/mediamon/collectors"
	"github.com/clambin/mediamon/collectors/bandwidth"
	"github.com/clambin/mediamon/collectors/connectivity"
	"github.com/clambin/mediamon/collectors/plex"
	"github.com/clambin/mediamon/collectors/transmission"
	"github.com/clambin/mediamon/collectors/xxxarr"
	"github.com/clambin/mediamon/collectors/xxxarr/scraper"
	scraperMock "github.com/clambin/mediamon/collectors/xxxarr/scraper/mocks"
	plex2 "github.com/clambin/mediamon/pkg/mediaclient/plex"
	plexMock "github.com/clambin/mediamon/pkg/mediaclient/plex/mocks"
	transmission2 "github.com/clambin/mediamon/pkg/mediaclient/transmission"
	transmissionMock "github.com/clambin/mediamon/pkg/mediaclient/transmission/mocks"
	"github.com/clambin/mediamon/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testCase struct {
	name   string
	config services.Config
}

var update = flag.Bool("update", false, "update .golden files")

var testCases = []testCase{
	{
		name: "full",
		config: services.Config{
			Transmission: transmission.Config{
				URL: "http://localhost",
			},
			Sonarr: xxxarr.Config{
				URL:    "http://localhost",
				APIKey: "1234",
			},
			Radarr: xxxarr.Config{
				URL:    "http://localhost",
				APIKey: "1234",
			},
			Plex: plex.Config{
				URL:      "http://localhost",
				UserName: "foo",
				Password: "bar",
			},
			OpenVPN: struct {
				Bandwidth    bandwidth.Config
				Connectivity connectivity.Config
			}{
				Bandwidth:    bandwidth.Config{FileName: "foo"},
				Connectivity: connectivity.Config{Proxy: "http://localhost", Token: "foo", Interval: time.Hour},
			},
		},
	},
	{
		name: "single",
		config: services.Config{
			Transmission: transmission.Config{
				URL: "http://localhost",
			},
		},
	},
}

func TestCreate(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := prometheus.NewRegistry()
			c := collectors.Create(&tt.config, r)
			assert.NotNil(t, c)
			buildUp(t, &c)

			gp := filepath.Join("testdata", strings.ToLower(t.Name())+".golden")
			if *update {
				metrics, err := r.Gather()
				require.NoError(t, err)

				buf := bytes.NewBuffer([]byte{})
				enc := expfmt.NewEncoder(buf, expfmt.FmtOpenMetrics)
				for _, entry := range metrics {
					err = enc.Encode(entry)
					require.NoError(t, err)
				}
				err = os.WriteFile(gp, buf.Bytes(), 0644)
				require.NoError(t, err)
			}

			f, err := os.Open(gp)
			require.NoError(t, err)
			err = testutil.GatherAndCompare(r, f)
			assert.NoError(t, err)

			//mock.AssertExpectationsForObjects(t, mocks...)
			tearDown(&c)
		})
	}
}

func buildUp(t *testing.T, c *collectors.Collectors) {
	if c.Transmission != nil {
		c.Transmission.API = createTransmissionMock(t)
	}
	if c.Sonarr != nil {
		c.Sonarr.Scraper = createSonarrMock(t)
	}
	if c.Radarr != nil {
		c.Radarr.Scraper = createRadarrMock(t)
	}
	if c.Plex != nil {
		c.Plex.API = createPlexMock(t)
	}
	if c.Bandwidth != nil {
		c.Bandwidth.Filename = createBandwidthFile()
	}
}

func tearDown(c *collectors.Collectors) {
	if c.Bandwidth != nil {
		_ = os.Remove(c.Bandwidth.Filename)
	}
}

func times() (n int) {
	n = 1
	if *update {
		n++
	}
	return n
}

func createTransmissionMock(t *testing.T) (m *transmissionMock.API) {
	m = transmissionMock.NewAPI(t)

	var o1 transmission2.SessionParameters
	o1.Arguments.Version = "foo"
	m.On("GetSessionParameters", mock.AnythingOfType("*context.emptyCtx")).Return(o1, nil).Times(times())

	var o2 transmission2.SessionStats
	o2.Arguments.ActiveTorrentCount = 1
	o2.Arguments.PausedTorrentCount = 2
	o2.Arguments.UploadSpeed = 25
	o2.Arguments.DownloadSpeed = 100
	m.On("GetSessionStatistics", mock.AnythingOfType("*context.emptyCtx")).Return(o2, nil).Times(times())

	return m
}

func createSonarrMock(t *testing.T) (m *scraperMock.Scraper) {
	m = scraperMock.NewScraper(t)
	m.On("Scrape", mock.AnythingOfType("*context.emptyCtx")).Return(scraper.Stats{
		URL:         "http://localhost",
		Version:     "foo",
		Monitored:   5,
		Unmonitored: 2,
	}, nil).Times(times())
	return m
}

func createRadarrMock(t *testing.T) (m *scraperMock.Scraper) {
	m = scraperMock.NewScraper(t)
	m.On("Scrape", mock.AnythingOfType("*context.emptyCtx")).Return(scraper.Stats{
		URL:         "http://localhost",
		Version:     "foo",
		Monitored:   2,
		Unmonitored: 5,
	}, nil).Times(times())
	return m
}

func createPlexMock(t *testing.T) (m *plexMock.API) {
	m = plexMock.NewAPI(t)
	m.On("GetIdentity", mock.AnythingOfType("*context.emptyCtx")).Return(plex2.IdentityResponse{}, nil).Times(times())
	m.On("GetSessions", mock.AnythingOfType("*context.emptyCtx")).Return(plex2.SessionsResponse{}, nil).Times(times())
	return m
}

func createBandwidthFile() string {
	f, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	_, _ = f.WriteString(`OpenVPN STATISTICS
Updated,Fri Dec 18 11:24:01 2020
TCP/UDP read bytes,5624951995
TCP/UDP write bytes,2048
END
`)
	_ = f.Close()
	return f.Name()
}
