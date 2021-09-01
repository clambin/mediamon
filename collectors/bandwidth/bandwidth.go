package bandwidth

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	filename string
}

type bandwidthStats struct {
	read    int64
	written int64
}

func NewCollector(filename string, _ time.Duration) prometheus.Collector {
	return &Collector{filename: filename}
}

func (coll *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- readMetric
	ch <- writeMetric
}

func (coll *Collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := coll.getStats()
	if err != nil {
		log.WithError(err).Warning("failed to collect bandwidth metrics")
		return
	}
	ch <- prometheus.MustNewConstMetric(readMetric, prometheus.GaugeValue, float64(stats.read))
	ch <- prometheus.MustNewConstMetric(writeMetric, prometheus.GaugeValue, float64(stats.written))
}

func (coll *Collector) getStats() (stats bandwidthStats, err error) {
	var file *os.File
	file, err = os.Open(coll.filename)

	if err != nil {
		return
	}

	fieldsFound := 0
	r := regexp.MustCompile(`^(.+),(\d+)$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, match := range r.FindAllStringSubmatch(scanner.Text(), -1) {
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
