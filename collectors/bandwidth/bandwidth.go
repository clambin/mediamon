package bandwidth

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	readMetric = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "tcp_udp_read_bytes_total"),
		"OpenVPN client bytes read",
		nil,
		nil,
	)
	writeMetric = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "tcp_udp_write_bytes_total"),
		"OpenVPN client bytes written",
		nil,
		nil,
	)
)

// Collector reads an openvpn status file and provides Prometheus metrics
type Collector struct {
	Filename string
	logger   *slog.Logger
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	FileName string
}

type bandwidthStats struct {
	read    int64
	written int64
}

// NewCollector creates a new Collector
func NewCollector(filename string) *Collector {
	return &Collector{
		Filename: filename,
		logger:   slog.Default().With("collector", "bandwidth"),
	}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- readMetric
	ch <- writeMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	stats, err := coll.getStats()
	if err != nil {
		// ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("mediamon_error", "Error getting bandwidth statistics", nil, nil), err)
		coll.logger.Error("failed to collect bandwidth metrics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
	coll.logger.Debug("bandwidth stats collected", "duration", time.Since(start))
}

var (
	statusFileRegEx = regexp.MustCompile(`^(.+),(\d+)$`)
)

func (coll *Collector) getStats() (bandwidthStats, error) {
	var stats bandwidthStats
	statusFile, err := os.Open(coll.Filename)
	if err != nil {
		return stats, err
	}
	defer func() { _ = statusFile.Close() }()

	fieldsFound := 0
	scanner := bufio.NewScanner(statusFile)
	for scanner.Scan() {
		line := scanner.Text()
		for _, match := range statusFileRegEx.FindAllStringSubmatch(line, -1) {
			value, _ := strconv.ParseInt(match[2], 10, 64)
			switch match[1] {
			case "TCP/UDP read bytes":
				stats.read = value
				fieldsFound++
			case "TCP/UDP write bytes":
				stats.written = value
				fieldsFound++
			}
		}
	}
	if fieldsFound != 2 {
		err = fmt.Errorf("not all fields were found in the openvpn status file")
	}

	return stats, err
}
