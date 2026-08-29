package bandwidth

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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
	logger   *slog.Logger
	Filename string
}

var _ prometheus.Collector = &Collector{}

// Config to create a Collector
type Config struct {
	FileName string
}

// NewCollector creates a new Collector
func NewCollector(filename string, logger *slog.Logger) prometheus.Collector {
	return &Collector{
		Filename: filename,
		logger:   logger,
	}
}

// Describe implements the prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- readMetric
	ch <- writeMetric
}

// Collect implements the prometheus.Collector interface
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	statusFile, err := os.Open(c.Filename)
	if err != nil {
		c.logger.Error("failed to open openvpn status file", "err", err)
		return
	}
	defer func() { _ = statusFile.Close() }()

	values, err := readClientStatusFile(statusFile)
	if err != nil {
		c.logger.Error("failed to read openvpn status file", "err", err)
		return
	}

	values.collect(ch)
}

var ignoredLines = map[string]struct{}{"OpenVPN STATISTICS": {}, "END": {}}

func readClientStatusFile(r io.Reader) (bandwidthStats, error) {
	values := make(bandwidthStats)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if _, ok := ignoredLines[line]; ok {
			continue
		}
		before, after, ok := strings.Cut(line, ",")
		if !ok {
			return nil, fmt.Errorf("invalid line %q", line)
		}
		if before == "Updated" {
			continue
		}
		value, err := strconv.ParseInt(after, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q: %w", after, err)
		}
		values[before] = value
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to read client status file: %w", err)
	}
	return values, nil
}

type bandwidthStats map[string]int64

func (s bandwidthStats) collect(ch chan<- prometheus.Metric) {
	const (
		bytesWritten = "TCP/UDP write bytes"
		bytesRead    = "TCP/UDP read bytes"
	)
	if value, ok := s[bytesWritten]; ok {
		ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(value))
	}
	if value, ok := s[bytesRead]; ok {
		ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(value))
	}
}
