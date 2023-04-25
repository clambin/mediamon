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
The  configuration file option specifies a yaml file to control tado-monitor's behaviour:

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
    # mediamon will connect to https://ipinfo.io through a proxy running inside the OpenVPN container
    # URL of the Proxy. If not set, connectivity won't be monitored
    proxy: <url>
    # Token supplied by ipinfo.io. You will need to register to obtain this
    token: <token>
    # interval limits how often connectivity is checked 
    interval: <duration>
```

If the filename is not specified on the command line, tado-monitor will look for a file `config.yaml` in the following directories:

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

```
mediamon_plex_session_count:                Active plex sessions by location (value indicates progress in the movie / episode)
mediamon_plex_transcoder_count:             Number of active transcode sessions
mediamon_plex_transcoder_speed:             Total speed of active transcoders
mediamon_plex_transcoder_count:             Number of transcode sessions
mediamon_plex_version:                      version info
mediamon_transmission_active_torrent_count: Number of active torrents
mediamon_transmission_download_speed:       Transmission download speed in bytes / sec
mediamon_transmission_paused_torrent_count: Number of paused torrents
mediamon_transmission_upload_speed:         Transmission upload speed in bytes / sec
mediamon_transmission_version:              version info
mediamon_xxxarr_calendar:                   Upcoming episodes / movies ("title" label contains the title)
mediamon_xxxarr_monitored_count:            Number of monitored series / movies
mediamon_xxxarr_queued_count:               Number of episodes / movies being downloaded
mediamon_xxxarr_queued_total_bytes:         Size of an episode / movie being downloaded ("title" label contains the title)
mediamon_xxxarr_queued_downloaded_bytes:    Downloaded size of an episode / movie being downloaded ("title" label contains the title)
mediamon_xxxarr_unmonitored_count:          Number of unmonitored series / movies
mediamon_xxxarr_version:                    version info
openvpn_client_status:                      OpenVPN client status
openvpn_client_tcp_udp_read_bytes_total:    OpenVPN client bytes read
openvpn_client_tcp_udp_write_bytes_total:   OpenVPN client bytes written
```

Additionally, the following metrics provide API usage statistics:

```
mediamon_request_duration_seconds Duration of API requests (summary metric)
mediamon_request_errors_total     API requests errors
```

### Grafana

[GitHub](https://github.com/clambin/mediamon/tree/master/assets/grafana/dashboards) contains a sample Grafana dashboard to visualize the scraped metrics.
Feel free to customize as you see fit.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
