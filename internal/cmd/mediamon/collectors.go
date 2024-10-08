package mediamon

import (
	"fmt"
	"github.com/clambin/mediamon/v2/internal/collectors/bandwidth"
	"github.com/clambin/mediamon/v2/internal/collectors/connectivity"
	"github.com/clambin/mediamon/v2/internal/collectors/plex"
	"github.com/clambin/mediamon/v2/internal/collectors/prowlarr"
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
		name: "transmission",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return transmission.NewCollector(url, logger)
		},
	},
	"sonarr.url": {
		name: "sonarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewSonarrCollector(url, v.GetString("sonarr.apikey"), logger)
		},
	},
	"radarr.url": {
		name: "radarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return xxxarr.NewRadarrCollector(url, v.GetString("radarr.apikey"), logger)
		},
	},
	"prowlarr.url": {
		name: "prowlarr",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return prowlarr.New(url, v.GetString("prowlarr.apikey"), logger)
		},
	},
	"plex.url": {
		name: "plex",
		make: func(url, version string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return plex.NewCollector(
				version,
				url,
				v.GetString("plex.username"),
				v.GetString("plex.password"),
				logger,
			), nil
		},
	},
	"openvpn.connectivity.proxy": {
		name: "vpn connectivity",
		make: func(url, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			proxy, err := parseProxy(url)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy. connectivity won't be monitored: %w", err)
			}
			return connectivity.NewCollector(
				proxy,
				v.GetDuration("openvpn.connectivity.interval"),
				logger,
			), nil
		},
	},
	"openvpn.bandwidth.filename": {
		name: "vpn bandwidth",
		make: func(target, _ string, v *viper.Viper, logger *slog.Logger) (prometheus.Collector, error) {
			return bandwidth.NewCollector(target, logger), nil
		},
	},
}

func createCollectors(version string, v *viper.Viper, logger *slog.Logger) []prometheus.Collector {
	collectors := make([]prometheus.Collector, 0, len(constructors))
	for key, c := range constructors {
		l := logger.With("collector", c.name)
		if value := v.GetString(key); value != "" {
			collector, err := c.make(value, version, v, l)
			if err != nil {
				l.Error("error creating collector", "err", err)
				continue
			}
			collectors = append(collectors, collector)
			l.Info("collector added", "source", value)
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
