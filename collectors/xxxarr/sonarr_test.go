package xxxarr_test

import (
	"github.com/clambin/mediamon/collectors/xxxarr"
	"testing"
	"time"
)

func TestSonarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewSonarrCollector("http://localhost:8888", "", 5*time.Minute)
	testCollectorDescribe(t, c, "constLabels: {application=\"sonarr\"}")
}

func TestSonarrCollector_Collect(t *testing.T) {
	c := xxxarr.NewSonarrCollector("", "", 5*time.Minute)
	c.(*xxxarr.SonarrCollector).Updater.API = &server{application: "sonarr"}
	testCollectorCollect(t, c, "sonarr")
}
