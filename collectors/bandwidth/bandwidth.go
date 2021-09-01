package bandwidth

import (
	"bufio"
	"github.com/clambin/mediamon/cache"
	"github.com/prometheus/client_golang/prometheus"
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

type Collector struct {
	cache.Cache
	filename string
}

type bandwidthStats struct {
	read    int64
	written int64
}

func NewCollector(filename string, interval time.Duration) prometheus.Collector {
	c := &Collector{filename: filename}
	c.Cache = *cache.New(interval, bandwidthStats{}, c.getStats)
	return c
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- readMetric
	ch <- writeMetric
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats := coll.Update().(bandwidthStats)

	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
}

func (coll *Collector) getStats() (interface{}, error) {
	var stats bandwidthStats
	file, err := os.Open(coll.filename)

	if err == nil {
		r := regexp.MustCompile(`^(.+),(\d+)$`)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			for _, match := range r.FindAllStringSubmatch(scanner.Text(), -1) {
				value, _ := strconv.ParseInt(match[2], 10, 64)
				switch match[1] {
				case "TCP/UDP read bytes":
					stats.read = value
				case "TCP/UDP write bytes":
					stats.written = value
				}
			}
		}
		_ = file.Close()
	}

	return stats, err
}
