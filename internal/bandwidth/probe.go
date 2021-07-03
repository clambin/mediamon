package bandwidth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/clambin/mediamon/internal/metrics"
)

// Probe to measure Plex metrics
type Probe struct {
	filename string
}

// NewProbe creates a new Probe
func NewProbe(filename string) *Probe {
	return &Probe{filename: filename}
}

// Run the probe. Collect all requires metrics
func (probe *Probe) Run(_ context.Context) error {
	var (
		err   error
		stats *openVPNStats
	)
	if stats, err = probe.getStats(); err == nil {
		metrics.OpenVPNClientReadTotal.Set(float64(stats.clientTCPUDPRead))
		metrics.OpenVPNClientWriteTotal.Set(float64(stats.clientTCPUDPWrite))
	} else {
		metrics.OpenVPNClientReadTotal.Set(0.0)
		metrics.OpenVPNClientWriteTotal.Set(0.0)
		log.WithField("err", err).Warning("Failed to get Bandwidth statistics")
	}

	return err
}

type openVPNStats struct {
	clientTCPUDPRead  int64
	clientTCPUDPWrite int64
}

func (stats *openVPNStats) String() string {
	return fmt.Sprintf("read=%d write=%d", stats.clientTCPUDPRead, stats.clientTCPUDPWrite)
}

func (probe *Probe) getStats() (*openVPNStats, error) {
	var stats = openVPNStats{}

	r := regexp.MustCompile(`^(.+),(\d+)$`)

	file, err := os.Open(probe.filename)
	if err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, match := range r.FindAllStringSubmatch(scanner.Text(), -1) {
				value, _ := strconv.ParseInt(match[2], 10, 64)
				switch match[1] {
				case "TCP/UDP read bytes":
					stats.clientTCPUDPRead = value
				case "TCP/UDP write bytes":
					stats.clientTCPUDPWrite = value
				}
			}
		}
	}

	log.WithFields(log.Fields{"err": err, "stats": stats.String()}).Debug("bandwidth getStats")

	return &stats, err
}
