package mediamon

import (
	"fmt"
	"github.com/clambin/mediamon/v2/internal/collectors/bandwidth"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/clambin/mediamon/v2/internal/collectors/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/transmission"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"log/slog"
	"net/url"
)

type constructor struct {
	name string
	make func(string, string, *viper.Viper, *slog.Logger) (prometheus.Collector, error)
}

var constructors = map[string]constructor{
	"transmission.url": {
		name: "Transmission",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return transmission.NewCollector(
				url,
				logger.With("collector", "transmission"),
			), nil
		},
	},
	"sonarr.url": {
		name: "Sonarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewSonarrCollector(
				url,
				v.GetString("sonarr.apikey"),
				logger.With("collector", "sonarr"),
			), nil
		},
	},
	"radarr.url": {
		name: "Radarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewRadarrCollector(
				url,
				v.GetString("radarr.apikey"),
				logger.With("collector", "radarr"),
			), nil
		},
	},
	"plex.url": {
		name: "Plex",
		make: func(url, version string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return plex.NewCollector(
				version,
				url,
				v.GetString("plex.username"),
				v.GetString("plex.password"),
				logger.With("collector", "plex"),
			), nil
		},
	},
	"openvpn.connectivity.proxy": {
		name: "VPN connectivity",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			proxy, err := parseProxy(url)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy. connectivity won't be monitored: %w", err)
			}
			return connectivity.NewCollector(
				v.GetString("openvpn.connectivity.token"),
				proxy,
				v.GetDuration("openvpn.connectivity.interval"),
				logger.With("collector", "connectivity"),
			), nil
		},
	},
	"openvpn.bandwidth.filename": {
		name: "VPN bandwidth",
		make: func(target, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return bandwidth.NewCollector(target, logger.With("collector", "bandwidth")), nil
		},
	},
}

func createCollectors(version string, v *viper.Viper, logger *slog.Logger) []prometheus.Collector {
	var collectors []prometheus.Collector

	for key, c := range constructors {
		l := logger.With("collector", c.name)
		if value := v.GetString(key); value != "" {
			collector, err := c.make(value, version, v, l)
			if err != nil {
				l.Error("error creating collector", "err", err)
				continue
			}
			l.Info("monitoring " + value)
			collectors = append(collectors, collector)
		}
	}
	return collectors
}

func parseProxy(proxyURL string) (*url.URL, error) {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	if proxy.Scheme == "" || proxy.Host == "" {
		return nil, fmt.Errorf("missing scheme / host")
	}
	return proxy, nil
}
