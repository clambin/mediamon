# mediamon
[![release](https://img.shields.io/github/v/tag/clambin/mediamon?color=green&label=release&style=plastic)](https://github.com/clambin/mediamon/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/mediamon?style=plastic)](https://app.codecov.io/gh/clambin/mediamon)
[![build](https://github.com/clambin/mediamon/workflows/Build/badge.svg)](https://github.com/clambin/mediamon/actions/workflows/build.yaml)
[![go report card](https://goreportcard.com/badge/github.com/clambin/mediamon/v2)](https://goreportcard.com/report/github.com/clambin/mediamon/v2)
[![license](https://img.shields.io/github/license/clambin/mediamon?style=plastic)](LICENSE.md)

Prometheus exporter for various media applications. Currently, supports Transmission, OpenVPN Client, Sonarr, Radarr, Prowlarr and Plex.

## Installation
Docker images are available on [ghcr.io](https://ghcr.io/clambin/mediamon).

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

prowlarr:
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
| mediamon_http_cache_hit_total | COUNTER | application, method, path|Number of times the cache was used |
| mediamon_http_cache_total | COUNTER | application, method, path|Number of times the cache was consulted |
| mediamon_http_request_duration_seconds | SUMMARY | application, code, method, path|duration of http requests |
| mediamon_http_requests_total | COUNTER | application, code, method, path|total number of http requests |
| mediamon_plex_library_bytes | GAUGE | library, url|Library size in bytes |
| mediamon_plex_library_count | GAUGE | library, url|Library size in number of entries |
| mediamon_plex_version | GAUGE | url, version|version info |
| mediamon_prowlarr_indexer_failed_grab_total | COUNTER | application, indexer, url|Total number of failed grabs from this indexer |
| mediamon_prowlarr_indexer_failed_query_total | COUNTER | application, indexer, url|Total number of failed queries to this indexer |
| mediamon_prowlarr_indexer_grab_total | COUNTER | application, indexer, url|Total number of grabs from this indexer |
| mediamon_prowlarr_indexer_query_total | COUNTER | application, indexer, url|Total number of queries to this indexer |
| mediamon_prowlarr_indexer_response_time | GAUGE | application, indexer, url|Average response time in seconds |
| mediamon_prowlarr_user_agent_grab_total | COUNTER | application, url, user_agent|Total number of grabs by user agent |
| mediamon_prowlarr_user_agent_query_total | COUNTER | application, url, user_agent|Total number of queries by user agent |
| mediamon_transmission_active_torrent_count | GAUGE | url|Number of active torrents |
| mediamon_transmission_download_speed | GAUGE | url|Transmission download speed in bytes / sec |
| mediamon_transmission_paused_torrent_count | GAUGE | url|Number of paused torrents |
| mediamon_transmission_upload_speed | GAUGE | url|Transmission upload speed in bytes / sec |
| mediamon_transmission_version | GAUGE | url, version|version info |
| mediamon_xxxarr_calendar | GAUGE | application, title, url|Upcoming episodes / movies |
| mediamon_xxxarr_health | GAUGE | application, type, url|Server health |
| mediamon_xxxarr_monitored_count | GAUGE | application, url|Number of Monitored series / movies |
| mediamon_xxxarr_queued_count | GAUGE | application, url|Episodes / movies being downloaded |
| mediamon_xxxarr_queued_downloaded_bytes | GAUGE | application, title, url|Downloaded size of episode / movie being downloaded in bytes |
| mediamon_xxxarr_queued_total_bytes | GAUGE | application, title, url|Size of episode / movie being downloaded in bytes |
| mediamon_xxxarr_unmonitored_count | GAUGE | application, url|Number of Unmonitored series / movies |
| mediamon_xxxarr_version | GAUGE | application, url, version|Version info |

### Grafana

[GitHub](https://github.com/clambin/mediamon/tree/master/assets/grafana/dashboards) contains a sample Grafana dashboard to visualize the scraped metrics.
Feel free to customize as you see fit.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
