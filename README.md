# github.com/clambin/mediamon
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/clambin/github.com/clambin/mediamon?color=green&label=Release&style=plastic)
![Codecov](https://img.shields.io/codecov/c/gh/clambin/github.com/clambin/mediamon?style=plastic)
![Build](https://github.com/clambin/github.com/clambin/mediamon/workflows/Build/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/clambin/github.com/clambin/mediamon)
![GitHub](https://img.shields.io/github/license/clambin/github.com/clambin/mediamon?style=plastic)

Prometheus exporter for various media applications.  Currently supports Transmission, OpenVPN Client, Sonarr, Radarr and Plex.

## Installation

A Docker image is available on [docker](https://hub.docker.com/r/clambin/github.com/clambin/mediamon).  Images are available for amd64 & arm32v7.

Alternatively, you can clone the repository from [github](https://github.com/clambin/github.com/clambin/mediamon) and build from source:

```
git clone https://github.com/clambin/github.com/clambin/mediamon.git
cd github.com/clambin/mediamon
go build
```

You will need to have go 1.15 installed on your system.

## Running github.com/clambin/mediamon
### Command-line options

The following command-line arguments can be passed:

```
$ ./github.com/clambin/mediamon
usage: github.com/clambin/mediamon --file=FILE [<flags>]

media monitor

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and --help-man).
  -v, --version    Show application version.
      --debug      Log debug messages
      --port=8080  API listener port
      --file=FILE  Service configuration file
```

### Configuration

The mandatory service configuration file configures which services github.com/clambin/mediamon should monitor:

```
transmission:
  # Transmission RPC URL, e.g. "http://192.168.0.1:9101/transmission/rpc"
  # If not set, Transmission won't be monitored
  url: <url>
  # How frequently to scrape transmission
  interval: <duration>

sonarr:
  # Sonarr URL. If not set, Sonarr won't be monitored
  url: <url>
  # How frequently to scrape Sonarr
  interval: <duration>
  # Sonarr API Key. See Sonarr / Settings / Security
  apikey: <key>

radarr:
  # All these are equivalent to sonarr
  url: <url>
  interval: <duration>
  apikey: <key>

plex:
  # Plex URL, e.g. http://192.168.0.11:32400 
  url: <url> 
  interval: <duration>
  # Your plex.tv user name and password
  username: <username>
  password: <password>

openvpn:
  # OpenVPN monitoring. Includes connectivity monitoring (up/down) and bandwidth consumption
  connectivity:
    # github.com/clambin/mediamon will connect to https://ipinfo.io through a proxy running inside the OpenVPN container
    # URL of the Proxy. If not set, connectivity won't be monitored
    proxy: <url>
    # Token supplied by ipinfo.io. You will need to register to obtain this
    token: <token>
    interval: <duration>
  bandwidth:
    # github.com/clambin/mediamon uses the OpenVPN status will to measure up/download bandwidth
    # filename contains the full path name of the client.status file. If not set, bandwidth won't be monitored
    filename: <file path>
    interval: <duration>>
```

### Prometheus

Add github.com/clambin/mediamon as a target to let Prometheus scrape the metrics into its database.
This highly depends on your particular Prometheus configuration. In it simplest form, add a new scrape target to `prometheus.yml`:

```
scrape_configs:
- job_name: github.com/clambin/mediamon
  static_configs:
  - targets: [ '<github.com/clambin/mediamon_host>:8080' ]
```


### Metrics

github.com/clambin/mediamon exposes the following metrics:

```
* mediaserver_active_torrent_count: Number of active torrents
* mediaserver_calendar_count: Number of upcoming episodes / movies
* mediaserver_download_speed: Transmission download speed in bytes / sec
* mediaserver_monitored_count: Number of monitored series / movies
* mediaserver_paused_torrent_count: Number of paused torrents
* mediaserver_plex_session_count: Active Plex sessions by user
* mediaserver_plex_transcoder_encoding_count: Number of active transcoders
* mediaserver_plex_transcoder_speed_total: Speed of active transcoders (total)
* mediaserver_plex_transcoder_type_count: Active Transcoder count by type
* mediaserver_queued_count: Number of queued torrents
* mediaserver_server_info: Server info. The 'version' label shows the current version of sonarr/radarr/transmission/plex
* mediaserver_unmonitored_count: Number of unmonitored series / movies
* mediaserver_upload_speed: Transmission upload speed in bytes / sec
* openvpn_client_status: OpenVPN Client Status
* openvpn_client_tcp_udp_read_bytes_total: OpenVPN client bytes read
* openvpn_client_tcp_udp_write_bytes_total: OpenVPN client bytes written
```

### Grafana

[Github](https://github.com/clambin/github.com/clambin/mediamon/tree/master/assets/grafana/dashboards) contains a sample Grafana dashboard to visualize the scraped metrics.
Feel free to customize as you see fit.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.