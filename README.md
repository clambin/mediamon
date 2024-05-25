# mediamon
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/mediamon?color=green&label=Release&style=plastic)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/mediamon?style=plastic)
![Build](https://github.com/clambin/mediamon/workflows/Build/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/mediamon)
![GitHub](https://img.shields.io/github/license/clambin/mediamon?style=plastic)

Prometheus exporter for various media applications. Currently, supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.

## Installation
Binaries are available on the [release](https://github.com/clambin/mediamon/releases) page. Docker images are available on [ghcr.io](https://ghcr.io/clambin/mediamon).

## Running mediamon
### Command-line options
The following command-line arguments can be passed:

```
Usage:
  mediamon [flags]

Flags:
      --config string   Configuration file
      --debug           Log debug messages
  -h, --help            help for mediamon
  -v, --version         version for mediamon
```

### Configuration
The  configuration file option specifies a yaml file to control mediamon's behaviour:

```
transmission:
  # Transmission RPC URL, e.g. "http://192.168.0.1:9101/transmission/rpc"
  # If not set, Transmission won't be monitored
  url: <url>

sonarr:
  # Sonarr URL. If not set, Sonarr won't be monitored
  url: <url>
  # Sonarr API Key. See Sonarr / Settings / Security
  apikey: <key>

radarr:
  # All these are equivalent to sonarr
  url: <url>
  apikey: <key>

plex:
  # Plex URL, e.g. http://192.168.0.11:32400 
  url: <url> 
  # Your plex.tv user name and password
  username: <username>
  password: <password>

openvpn:
  bandwidth:
    # mediamon uses the OpenVPN status will to measure up/download bandwidth
    # filename contains the full path name of the client.status file. If not set, bandwidth won't be monitored
    filename: <file path>
  # OpenVPN monitoring. Includes connectivity monitoring (up/down) and bandwidth consumption
  connectivity:
    # mediamon will connect to http://ip-api.com through a proxy running inside the OpenVPN container
    # URL of the Proxy. If not set, connectivity won't be monitored
    proxy: <url>
    # interval limits how often connectivity is checked 
    interval: <duration>
```

If the filename is not specified on the command line, mediamon will look for a file `config.yaml` in the following directories:

```
/etc/mediamon
$HOME/.mediamon
.
```

Any value in the configuration file may be overriden by setting an environment variable with a prefix `MEDIAMON_`. 
E.g. to avoid setting your Sonarr API key in the configuration file, set the following environment variables:

```
export MEDIAMON_SONAR.APIKEY="your-sonarr-apikey"
```

### Prometheus
Add mediamon as a target to let Prometheus scrape the metrics into its database.
This highly depends on your particular Prometheus configuration. In its simplest form, add a new scrape target to `prometheus.yml`:

```
scrape_configs:
- job_name: mediamon
  static_configs:
  - targets: [ '<mediamon_host>:8080' ]
```


### Metrics

mediamon exposes the following metrics:

| metric | type |  labels | help |
| --- | --- |  --- | --- |
| mediamon_api_cache_hit_total | COUNTER | application, method, path|Number of times the cache was used |
| mediamon_api_cache_total | COUNTER | application, method, path|Number of times the cache was consulted |
| mediamon_api_errors_total | COUNTER | application, method, path|Number of failed HTTP calls |
| mediamon_api_latency | SUMMARY | application, method, path|latency of HTTP calls |
| mediamon_plex_library_bytes | GAUGE | library, url|Library size in bytes |
| mediamon_plex_library_count | GAUGE | library, url|Library size in number of entries |
| mediamon_plex_session_bandwidth | GAUGE | address, audioCodec, lat, location, lon, mode, player, title, url, user, videoCodec|Active Plex session Bandwidth usage (in kbps) |
| mediamon_plex_session_count | GAUGE | address, audioCodec, lat, location, lon, mode, player, title, url, user, videoCodec|Active Plex session progress |
| mediamon_plex_version | GAUGE | url, version|version info |
| mediamon_transmission_active_torrent_count | GAUGE | url|Number of active torrents |
| mediamon_transmission_download_speed | GAUGE | url|Transmission download speed in bytes / sec |
| mediamon_transmission_paused_torrent_count | GAUGE | url|Number of paused torrents |
| mediamon_transmission_upload_speed | GAUGE | url|Transmission upload speed in bytes / sec |
| mediamon_transmission_version | GAUGE | url, version|version info |
| mediamon_xxxarr_monitored_count | GAUGE | application, url|Number of Monitored series / movies |
| mediamon_xxxarr_queued_count | GAUGE | application, url|Episodes / movies being downloaded |
| mediamon_xxxarr_unmonitored_count | GAUGE | application, url|Number of Unmonitored series / movies |
| mediamon_xxxarr_version | GAUGE | application, url, version|Version info |

### Grafana

[GitHub](https://github.com/clambin/mediamon/tree/master/assets/grafana/dashboards) contains a sample Grafana dashboard to visualize the scraped metrics.
Feel free to customize as you see fit.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
