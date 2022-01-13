package xxxarr_test

import (
	"github.com/clambin/mediamon/collectors/xxxarr"
	"testing"
	"time"
)

func TestRadarrCollector_Describe(t *testing.T) {
	c := xxxarr.NewRadarrCollector("http://localhost:8888", "", 5*time.Minute)
	testCollectorDescribe(t, c, "constLabels: {application=\"radarr\"}")
}

func TestRadarrCollector_Collect(t *testing.T) {
	c := xxxarr.NewRadarrCollector("", "", 5*time.Minute)
	c.(*xxxarr.RadarrCollector).Updater.API = &server{application: "sonarr"}

	testCollectorCollect(t, c, "radarr")
}
