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
	make func(string, string, *viper.Viper, *slog.Logger) prometheus.Collector
}

var constructors = map[string]constructor{
	"transmission.url": {
		name: "Transmission",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			return transmission.NewCollector(
				url,
				logger.With("collector", "transmission"),
			)
		},
	},
	"sonarr.url": {
		name: "Sonarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			return xxxarr.NewSonarrCollector(
				url,
				v.GetString("sonarr.apikey"),
				logger.With("collector", "sonarr"),
			)
		},
	},
	"radarr.url": {
		name: "Radarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			return xxxarr.NewRadarrCollector(
				url,
				v.GetString("radarr.apikey"),
				logger.With("collector", "radarr"),
			)
		},
	},
	"plex.url": {
		name: "Plex",
		make: func(url, version string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			return plex.NewCollector(
				version,
				url,
				v.GetString("plex.username"),
				v.GetString("plex.password"),
				logger.With("collector", "plex"),
			)
		},
	},
	"openvpn.connectivity.proxy": {
		name: "VPN connectivity",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			proxy, err := parseProxy(url)
			if err != nil {
				logger.Error("invalid proxy. connectivity won't be monitored", "err", err)
				return nil
			}
			return connectivity.NewCollector(
				v.GetString("openvpn.connectivity.token"),
				proxy,
				v.GetDuration("openvpn.connectivity.interval"),
				logger.With("collector", "connectivity"),
			)
		},
	},
	"openvpn.bandwidth.filename": {
		name: "VPN bandwidth",
		make: func(target, _ string, v *viper.Viper, logger *slog.Logger) prometheus.Collector {
			return bandwidth.NewCollector(target, logger.With("collector", "bandwidth"))
		},
	},
}

func createCollectors(version string, v *viper.Viper, logger *slog.Logger) []prometheus.Collector {
	var collectors []prometheus.Collector

	for key, c := range constructors {
		if value := v.GetString(key); value != "" {
			logger.Info("monitoring "+key, "target", key)
			if collector := c.make(value, version, v, logger.With("collector", c.name)); collector != nil {
				collectors = append(collectors, collector)
			}
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
