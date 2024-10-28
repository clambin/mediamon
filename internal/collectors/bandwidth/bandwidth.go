package bandwidth

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
func NewCollector(filename string, logger *slog.Logger) *Collector {
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
	var stats bandwidthStats
	if err == nil {
		stats, err = readStats(statusFile)
		_ = statusFile.Close()
	}
	if err != nil {
		c.logger.Error("failed to collect bandwidth metrics", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
}

func readStats(r io.Reader) (bandwidthStats, error) {
	values, err := readClientStatusFile(r)
	if err != nil {
		return bandwidthStats{}, err
	}
	var stats bandwidthStats
	var ok bool
	if stats.written, ok = values["TCP/UDP write bytes"]; !ok {
		return bandwidthStats{}, errors.New("TCP/UDP write bytes not found")
	}
	if stats.read, ok = values["TCP/UDP read bytes"]; !ok {
		return bandwidthStats{}, errors.New("TCP/UDP read bytes not found")
	}
	return stats, nil
}

var ignoredLines = map[string]struct{}{"OpenVPN STATISTICS": {}, "END": {}}

func readClientStatusFile(r io.Reader) (map[string]int64, error) {
	values := make(map[string]int64)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if _, ok := ignoredLines[line]; ok {
			continue
		}
		idx := strings.IndexByte(line, ',')
		if idx == -1 {
			return nil, fmt.Errorf("invalid line %q", line)
		}
		if line[:idx] == "Updated" {
			continue
		}
		value, err := strconv.ParseInt(line[idx+1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q: %w", line[idx+1:], err)
		}
		values[line[:idx]] = value
	}
	return values, nil
}
