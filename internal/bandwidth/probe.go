package bandwidth

import (
	"bufio"
	"os"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"mediamon/internal/metrics"
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
func (probe *Probe) Run() {
	if stats, err := probe.getStats(); err == nil {
		metrics.OpenVPNClientReadTotal.Set(float64(stats.clientTcpUdpRead))
		metrics.OpenVPNClientWriteTotal.Set(float64(stats.clientTcpUdpWrite))
	} else {
		metrics.OpenVPNClientReadTotal.Set(0.0)
		metrics.OpenVPNClientWriteTotal.Set(0.0)
		log.Warningf("%s", err.Error())
	}
}

type openVPNStats struct {
	clientTcpUdpRead  int64
	clientTcpUdpWrite int64
}

func (probe *Probe) getStats() (openVPNStats, error) {
	var stats = openVPNStats{}

	r := regexp.MustCompile(`^(.+),(\d+)$`)

	file, err := os.Open(probe.filename)
	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, match := range r.FindAllStringSubmatch(scanner.Text(), -1) {
				value, _ := strconv.ParseInt(match[2], 10, 64)
				switch match[1] {
				case "TCP/UDP read bytes":
					// log.Debugf("clientTcpUdpRead: %s -> %s -> %d", scanner.Text(), match[2], value)
					stats.clientTcpUdpRead = value
				case "TCP/UDP write bytes":
					// log.Debugf("clientTcpUdpWrite: %s -> %s -> %d", scanner.Text(), match[2], value)
					stats.clientTcpUdpWrite = value
				}
			}
		}
	}

	return stats, err
}
