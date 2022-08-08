package bandwidth

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	return &Collector{Filename: filename}
}

// Describe implements the prometheus.Collector interface
func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- readMetric
	ch <- writeMetric
}

// Collect implements the prometheus.Collector interface
func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := coll.getStats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(
			prometheus.NewDesc("mediamon_error",
				"Error getting bandwidth statistics", nil, nil),
			err)
		log.WithError(err).Warning("failed to collect bandwidth metrics")
		return
	}
	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
}

func (coll *Collector) getStats() (stats bandwidthStats, err error) {
	var file *os.File
	file, err = os.Open(coll.Filename)

	if err != nil {
		return
	}

	fieldsFound := 0
	r := regexp.MustCompile(`^(.+),(\d+)$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, match := range r.FindAllStringSubmatch(line, -1) {
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
	_ = file.Close()

	if fieldsFound != 2 {
		err = fmt.Errorf("not all fields were found in the openvpn status file")
	}

	return
}
