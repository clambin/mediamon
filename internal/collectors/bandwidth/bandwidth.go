package bandwidth

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
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
	stats, err := c.getStats(c.Filename)
	if err != nil {
		c.logger.Error("failed to collect bandwidth metrics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
}

var (
	statusFileRegEx = regexp.MustCompile(`^(.+),(\d+)$`)
)

func (c *Collector) getStats(filename string) (bandwidthStats, error) {
	statusFile, err := os.Open(filename)
	if err != nil {
		return bandwidthStats{}, err
	}
	defer func() { _ = statusFile.Close() }()

	return c.readStats(statusFile)
}

func (c *Collector) readStats(r io.Reader) (stats bandwidthStats, err error) {
	var fieldsFound int
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		for _, match := range statusFileRegEx.FindAllStringSubmatch(scanner.Text(), -1) {
			// regexp only matches on valid numbers, so ParseInt won't fail
			value, _ := strconv.ParseInt(match[2], 10, 64)
			switch match[1] {
			case "TCP/UDP read bytes":
				stats.read = value
				fieldsFound |= 0x1
			case "TCP/UDP write bytes":
				stats.written = value
				fieldsFound |= 0x2
			}
		}
	}
	if fieldsFound != 0x3 {
		err = fmt.Errorf("not all fields were found in the openvpn status file")
	}
	return stats, err
}
