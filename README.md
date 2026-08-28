# mediamon

[![release](https://img.shields.io/github/v/tag/clambin/mediamon?color=green&label=release&style=plastic)](https://github.com/clambin/mediamon/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/mediamon?style=plastic)](https://app.codecov.io/gh/clambin/mediamon)
[![build](https://github.com/clambin/mediamon/workflows/Build/badge.svg)](https://github.com/clambin/mediamon/actions/workflows/build.yaml)
[![license](https://img.shields.io/github/license/clambin/mediamon?style=plastic)](LICENSE.md)

Prometheus exporter for various media applications. Currently, it supports Transmission, OpenVPN Client, Sonarr, Radarr,
Prowlarr and Plex.

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

The configuration file option specifies a YAML file to control mediamon's behavior:

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
  # Your PMS Token. See https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token.
  # If you don't want to use your token here, leave this blank and provide your Plex username & password instead.
  # Optionally, enable JWT authentication to use your Plex username & password only on first login.
  token: <token> 
  # Plex username. Used to authenticate with Plex
  username: my-user
  # Plex username. Used to authenticate with Plex
  password: my-password
  # Plex Client Identifier. Used to authenticate with Plex.
  # When using JWT authentication, you *must* provide a client-id to prevent mediamon of logging in each time it starts.
  client-id: 2194c117-e4eb-4223-af4e-f924e2234d21
  # Plex JWT authentication. This is the new authentication method recommended by Plex.
  # The advantage of JWT authentication is that it only uses your Plex username & password on first login.
  # After the first login, mediamon will store a JWT token and use that to authenticate with Plex.
  jwt:
    # Enable JWT authentication.
    enable: true
    # Path to store the JWT token (encrypted).
    path: "storage.enc"
    # Passphrase to encrypt the JWT token.
    passphrase: "my-very-insecure-passphrase"
openvpn:
  bandwidth:
    # mediamon uses the OpenVPN status will to measure up/download bandwidth
    # filename contains the full path name of the client.status file. If not set, bandwidth won't be monitored
    filename: <file path>
  # OpenVPN monitoring. Includes connectivity monitoring (up/down) and bandwidth consumption
  connectivity:
    # mediamon will connect to https://ip-api.com through a proxy running inside the OpenVPN container
    # URL of the Proxy. If not set, connectivity won't be monitored
    proxy: <url>
    # interval limits how often connectivity is checked 
    interval: <duration>
```

If the filename is not specified on the command line, mediamon will look for a file `config.yaml` in the following
directories:

```
/etc/mediamon
$HOME/.mediamon
.
```

Any value in the configuration file may be overridden by setting an environment variable with a prefix `MEDIAMON_`. 
E.g., to avoid setting your Sonarr API key in the configuration file, set the following environment variables:

```
export MEDIAMON_SONAR.APIKEY="your-sonarr-apikey"
```

### Prometheus

Add mediamon as a target to let Prometheus scrape the metrics into its database. This highly depends on your particular
Prometheus configuration. In its simplest form, add a new scrape target to `prometheus.yml`:

```
scrape_configs:
- job_name: mediamon
  static_configs:
  - targets: [ '<mediamon_host>:8080' ]
```

### Metrics

mediamon exposes the following metrics:

| metric                                       | type    | labels                       | help                                           |
|----------------------------------------------|---------|------------------------------|------------------------------------------------|
| mediamon_http_request_duration_seconds       | SUMMARY | application, code, method    | duration of http requests                      |
| mediamon_http_requests_total                 | COUNTER | application, code, method    | total number of http requests                  |
| mediamon_plex_episode_count                  | GAUGE   | library, url                 | Total number of episodes in Plex library       |
| mediamon_plex_library_bytes                  | GAUGE   | library, url                 | Library size in bytes                          |
| mediamon_plex_movie_count                    | GAUGE   | library, url                 | Total number of movies in Plex library         |
| mediamon_plex_show_count                     | GAUGE   | library, url                 | Total number of shows in Plex library          |
| mediamon_plex_version                        | GAUGE   | url, version                 | version info                                   |
| mediamon_prowlarr_indexer_failed_grab_total  | COUNTER | application, indexer, url    | Total number of failed grabs from this indexer |
| mediamon_prowlarr_indexer_failed_query_total | COUNTER | application, indexer, url    | Total number of failed queries to this indexer |
| mediamon_prowlarr_indexer_grab_total         | COUNTER | application, indexer, url    | Total number of grabs from this indexer        |
| mediamon_prowlarr_indexer_query_total        | COUNTER | application, indexer, url    | Total number of queries to this indexer        |
| mediamon_prowlarr_indexer_response_time      | GAUGE   | application, indexer, url    | Average response time in seconds               |
| mediamon_prowlarr_user_agent_grab_total      | COUNTER | application, url, user_agent | Total number of grabs by user agent            |
| mediamon_prowlarr_user_agent_query_total     | COUNTER | application, url, user_agent | Total number of queries by user agent          |
| mediamon_transmission_active_torrent_count   | GAUGE   | url                          | Number of active torrents                      |
| mediamon_transmission_download_speed         | GAUGE   | url                          | Transmission download speed in bytes / sec     |
| mediamon_transmission_paused_torrent_count   | GAUGE   | url                          | Number of paused torrents                      |
| mediamon_transmission_upload_speed           | GAUGE   | url                          | Transmission upload speed in bytes / sec       |
| mediamon_transmission_version                | GAUGE   | url, version                 | version info                                   |
| mediamon_xxxarr_calendar                     | GAUGE   | application, title, url      | Upcoming episodes / movies                     |
| mediamon_xxxarr_health                       | GAUGE   | application, type, url       | Server health                                  |
| mediamon_xxxarr_monitored_count              | GAUGE   | application, url             | Number of Monitored series / movies            |
| mediamon_xxxarr_queued_count                 | GAUGE   | application, url             | Episodes / movies being downloaded             |
| mediamon_xxxarr_unmonitored_count            | GAUGE   | application, url             | Number of Unmonitored series / movies          |
| mediamon_xxxarr_version                      | GAUGE   | application, url, version    | Version info                                   |

### Grafana

[GitHub](https://github.com/clambin/mediamon/tree/master/assets/grafana/dashboards) contains a sample Grafana dashboard
to visualize the scraped metrics. Feel free to customize as you see fit.

## Authors

* **Christophe Lambin**

## License

This project is licensed under the MIT License – see the [LICENSE.md](LICENSE.md) file for details.
