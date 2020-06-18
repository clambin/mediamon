from src.mediacentre import TransmissionProbe, MonitorProbe


class FakeResponse:
    def __init__(self, status_code, headers, text):
        self.status_code = status_code
        self.headers = headers
        self.text = text

    def json(self):
        return self.text


class TransmissionTestProbe(TransmissionProbe):
    def __init__(self, host):
        super().__init__(host)

    def post(self, endpoint=None, headers=None, body=None):
        if 'X-Transmission-Session-Id' not in headers or headers['X-Transmission-Session-Id'] == '':
            return FakeResponse(409, {'X-Transmission-Session-Id': 'NewKey'}, {})
        else:
            return FakeResponse(200, {}, {'arguments': {
                "activeTorrentCount": 1,
                "cumulative-stats": {
                    "downloadedBytes": 259842832295,
                    "filesAdded": 218,
                    "secondsActive": 3106063,
                    "sessionCount": 19,
                    "uploadedBytes": 67534137454
                },
                "current-stats": {
                    "downloadedBytes": 53505238629,
                    "filesAdded": 39,
                    "secondsActive": 508967,
                    "sessionCount": 1,
                    "uploadedBytes": 14868574785
                },
                "downloadSpeed": 1000,
                "pausedTorrentCount": 2,
                "torrentCount": 3,
                "uploadSpeed": 500
            }})


class MonitorTestProbe(MonitorProbe):
    def __init__(self, host, name, api_key):
        super().__init__(host, name, api_key)

    def get(self, endpoint=None, headers=None, body=None, params=None):
        if self.app == self.App.sonarr:
            if endpoint == 'api/calendar':
                return FakeResponse(200, {}, [{
                    "seriesId": 33,
                    "episodeFileId": 3923,
                    "seasonNumber": 7,
                    "episodeNumber": 12,
                    "title": "XXXX",
                    "airDate": "2020-04-16",
                    "airDateUtc": "2020-04-17T00:30:00Z",
                    "overview": "XXXX",
                    "episodeFile": {
                        "seriesId": 33,
                        "seasonNumber": 7,
                        "relativePath": "XXXX",
                        "path": "XXXX",
                        "size": 546775439,
                        "dateAdded": "2020-04-17T03:50:58.417555Z",
                        "sceneName": "XXXX",
                        "quality": {
                            "quality": {
                                "id": 4,
                                "name": "HDTV-720p",
                                "source": "television",
                                "resolution": 720
                            },
                            "revision": {
                                "version": 1,
                                "real": 0
                            }
                        },
                        "mediaInfo": {
                            "audioChannels": 5.1,
                            "audioCodec": "AC3",
                            "videoCodec": "x264"
                        },
                        "originalFilePath": "XXXX",
                        "qualityCutoffNotMet": False,
                        "id": 3923
                    },
                    "hasFile": False,
                    "monitored": True,
                    "absoluteEpisodeNumber": 142,
                    "unverifiedSceneNumbering": False,
                    "series": {
                        "title": "XXXX",
                        "sortTitle": "XXXX",
                        "seasonCount": 7,
                        "status": "continuing",
                        "overview": "XXXX",
                        "network": "XXXX",
                        "airTime": "20:30",
                        "images": [
                            {
                                "coverType": "fanart",
                                "url": "XXXX",
                            },
                            {
                                "coverType": "banner",
                                "url": "XXXX",
                            },
                            {
                                "coverType": "poster",
                                "url": "XXXX",
                            }
                        ],
                        "seasons": [
                            {
                                "seasonNumber": 0,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 1,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 2,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 3,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 4,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 5,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 6,
                                "monitored": False
                            },
                            {
                                "seasonNumber": 7,
                                "monitored": True
                            }
                        ],
                        "year": 2013,
                        "path": "XXXX",
                        "profileId": 6,
                        "seasonFolder": True,
                        "monitored": True,
                        "useSceneNumbering": False,
                        "runtime": 25,
                        "tvdbId": 1,
                        "tvRageId": 1,
                        "tvMazeId": 1,
                        "firstAired": "2013-08-10T22:00:00Z",
                        "lastInfoSync": "2020-04-17T12:20:35.667424Z",
                        "seriesType": "standard",
                        "cleanTitle": "XXXX",
                        "imdbId": "XXXX",
                        "titleSlug": "XXXX",
                        "certification": "TV-14",
                        "genres": [
                            "Action",
                            "Comedy",
                            "Crime"
                        ],
                        "tags": [],
                        "added": "2020-02-01T20:27:00.262442Z",
                        "ratings": {
                            "votes": 2679,
                            "value": 8.6
                        },
                        "qualityProfileId": 6,
                        "id": 1
                    },
                    "id": 1
                }])
            elif endpoint == 'api/queue':
                return FakeResponse(200, {}, [{
                    "series": {
                        "title": "XXXX",
                        "sortTitle": "XXXX",
                        "seasonCount": 1,
                        "status": "continuing",
                        "overview": "XXXX",
                        "network": "XXXX",
                        "airTime": "22:30",
                        "images": [
                            {
                                "coverType": "fanart",
                                "url": "XXXX"
                            },
                            {
                                "coverType": "banner",
                                "url": "XXXX"
                            },
                            {
                                "coverType": "poster",
                                "url": "XXXX",
                            }
                        ],
                        "seasons": [
                            {
                                "seasonNumber": 1,
                                "monitored": True
                            }
                        ],
                        "year": 2020,
                        "path": "XXXX",
                        "profileId": 3,
                        "seasonFolder": True,
                        "monitored": True,
                        "useSceneNumbering": False,
                        "runtime": 30,
                        "tvdbId": 1,
                        "tvRageId": 0,
                        "tvMazeId": 1,
                        "firstAired": "2020-04-11T22:00:00Z",
                        "lastInfoSync": "2020-04-18T21:26:31.930624Z",
                        "seriesType": "standard",
                        "cleanTitle": "XXXX",
                        "imdbId": "XXXX",
                        "titleSlug": "XXXX",
                        "certification": "XXXX",
                        "genres": [
                            "Comedy",
                            "Drama",
                            "Romance",
                            "Thriller"
                        ],
                        "tags": [],
                        "added": "2020-04-18T21:26:31.622049Z",
                        "ratings": {
                            "votes": 0,
                            "value": 0.0
                        },
                        "qualityProfileId": 3,
                        "id": 40
                    },
                    "episode": {
                        "seriesId": 40,
                        "episodeFileId": 0,
                        "seasonNumber": 1,
                        "episodeNumber": 1,
                        "title": "XXXX",
                        "airDate": "2020-04-12",
                        "airDateUtc": "2020-04-13T02:30:00Z",
                        "overview": "XXXX",
                        "hasFile": False,
                        "monitored": True,
                        "unverifiedSceneNumbering": False,
                        "lastSearchTime": "2020-04-18T21:26:32.814329Z",
                        "id": 1736
                    },
                    "quality": {
                        "quality": {
                            "id": 5,
                            "name": "WEBDL-720p",
                            "source": "web",
                            "resolution": 720
                        },
                        "revision": {
                            "version": 2,
                            "real": 0
                        }
                    },
                    "size": 489530639.0,
                    "title": "XXXX",
                    "sizeleft": 485965824.0,
                    "timeleft": "03:47:14",
                    "estimatedCompletionTime": "2020-04-19T01:15:20.195938Z",
                    "status": "Downloading",
                    "trackedDownloadStatus": "Ok",
                    "statusMessages": [],
                    "downloadId": "XXXX",
                    "protocol": "torrent",
                    "id": 1
                }])
            elif endpoint == 'api/series':
                return FakeResponse(200, {}, [{
                    "title": "XXXX",
                    "alternateTitles": [],
                    "sortTitle": "XXXX",
                    "seasonCount": 2,
                    "totalEpisodeCount": 17,
                    "episodeCount": 16,
                    "episodeFileCount": 16,
                    "sizeOnDisk": 29496630349,
                    "status": "continuing",
                    "overview": "XXXX",
                    "previousAiring": "2019-04-29T01:00:00Z",
                    "network": "XXXX",
                    "airTime": "21:00",
                    "images": [
                        {
                            "coverType": "fanart",
                            "url": "XXXX",
                        },
                        {
                            "coverType": "banner",
                            "url": "XXXX",
                        },
                        {
                            "coverType": "poster",
                            "url": "XXXX",
                        }
                    ],
                    "seasons": [
                        {
                            "seasonNumber": 0,
                            "monitored": False,
                            "statistics": {
                                "episodeFileCount": 0,
                                "episodeCount": 0,
                                "totalEpisodeCount": 1,
                                "sizeOnDisk": 0,
                                "percentOfEpisodes": 0.0
                            }
                        },
                        {
                            "seasonNumber": 1,
                            "monitored": False,
                            "statistics": {
                                "previousAiring": "2017-06-19T01:00:00Z",
                                "episodeFileCount": 8,
                                "episodeCount": 8,
                                "totalEpisodeCount": 8,
                                "sizeOnDisk": 16443635433,
                                "percentOfEpisodes": 100.0
                            }
                        },
                        {
                            "seasonNumber": 2,
                            "monitored": False,
                            "statistics": {
                                "previousAiring": "2019-04-29T01:00:00Z",
                                "episodeFileCount": 8,
                                "episodeCount": 8,
                                "totalEpisodeCount": 8,
                                "sizeOnDisk": 13052994916,
                                "percentOfEpisodes": 100.0
                            }
                        }
                    ],
                    "year": 2017,
                    "path": "XXXX",
                    "profileId": 3,
                    "seasonFolder": True,
                    "monitored": False,
                    "useSceneNumbering": False,
                    "runtime": 55,
                    "tvdbId": 1,
                    "tvRageId": 1,
                    "tvMazeId": 1,
                    "firstAired": "2017-04-29T22:00:00Z",
                    "lastInfoSync": "2020-04-17T12:20:33.392966Z",
                    "seriesType": "standard",
                    "cleanTitle": "XXXX",
                    "imdbId": "XXXX",
                    "titleSlug": "XXXX",
                    "certification": "XXXX",
                    "genres": [
                        "Action",
                        "Adventure",
                        "Fantasy",
                        "Suspense"
                    ],
                    "tags": [],
                    "added": "2017-04-29T21:26:52.21431Z",
                    "ratings": {
                        "votes": 796,
                        "value": 8.3
                    },
                    "qualityProfileId": 3,
                    "id": 6
                }, {
                    "title": "XXXX",
                    "alternateTitles": [
                        {
                            "title": "XXXX",
                            "seasonNumber": -1
                        }, {
                            "title": "XXXX",
                            "seasonNumber": -1
                        }, {
                            "title": "XXXX",
                            "seasonNumber": -1
                        }
                    ],
                    "sortTitle": "XXXX",
                    "seasonCount": 3,
                    "totalEpisodeCount": 62,
                    "episodeCount": 13,
                    "episodeFileCount": 13,
                    "sizeOnDisk": 7505643639,
                    "status": "continuing",
                    "overview": "XXXX",
                    "previousAiring": "2019-08-14T07:00:00Z",
                    "network": "XXXX",
                    "airTime": "03:00",
                    "images": [{
                        "coverType": "fanart",
                        "url": "XXXX"
                    }, {
                        "coverType": "banner",
                        "url": "XXXX"
                    }, {
                        "coverType": "poster",
                        "url": "XXXX"
                    }],
                    "seasons": [{
                        "seasonNumber": 0,
                        "monitored": False,
                        "statistics": {
                            "episodeFileCount": 0,
                            "episodeCount": 0,
                            "totalEpisodeCount": 26,
                            "sizeOnDisk": 0,
                            "percentOfEpisodes": 0.0
                        }
                    }, {
                        "seasonNumber": 1,
                        "monitored": False,
                        "statistics": {
                            "episodeFileCount": 0,
                            "episodeCount": 0,
                            "totalEpisodeCount": 10,
                            "sizeOnDisk": 0,
                            "percentOfEpisodes": 0.0
                        }
                    }, {
                        "seasonNumber": 2,
                        "monitored": False,
                        "statistics": {
                            "episodeFileCount": 0,
                            "episodeCount": 0,
                            "totalEpisodeCount": 13,
                            "sizeOnDisk": 0,
                            "percentOfEpisodes": 0.0
                        }
                    }, {
                        "seasonNumber": 3,
                        "monitored": True,
                        "statistics": {
                            "previousAiring": "2019-08-14T07:00:00Z",
                            "episodeFileCount": 13,
                            "episodeCount": 13,
                            "totalEpisodeCount": 13,
                            "sizeOnDisk": 7505643639,
                            "percentOfEpisodes": 100.0
                        }
                    }
                    ],
                    "year": 2017,
                    "path": "XXXX",
                    "profileId": 3,
                    "seasonFolder": True,
                    "monitored": True,
                    "useSceneNumbering": False,
                    "runtime": 55,
                    "tvdbId": 1,
                    "tvRageId": 1,
                    "tvMazeId": 1,
                    "firstAired": "2017-04-25T22:00:00Z",
                    "lastInfoSync": "2020-04-17T12:20:37.020309Z",
                    "seriesType": "standard",
                    "cleanTitle": "XXXX",
                    "imdbId": "XXXX",
                    "titleSlug": "XXXX",
                    "certification": "XXXX",
                    "genres": ["Drama"],
                    "tags": [],
                    "added": "2017-07-14T08:58:45.103329Z",
                    "ratings": {
                        "votes": 891,
                        "value": 8.9
                    },
                    "qualityProfileId": 3,
                    "id": 9
                }])
        elif self.app == self.App.radarr:
            if endpoint == 'api/calendar':
                return FakeResponse(200, {}, [])
            elif endpoint == 'api/queue':
                return FakeResponse(200, {}, [
                    {
                        "movie": {
                            "title": "XXXX",
                            "alternativeTitles": [
                                {
                                    "sourceType": "tmdb",
                                    "movieId": 1,
                                    "title": "XXXX",
                                    "sourceId": 1,
                                    "votes": 0,
                                    "voteCount": 0,
                                    "language": "english",
                                    "id": 1
                                },
                                {
                                    "sourceType": "tmdb",
                                    "movieId": 1,
                                    "title": "XXXX",
                                    "sourceId": 1,
                                    "votes": 0,
                                    "voteCount": 0,
                                    "language": "english",
                                    "id": 1
                                },
                                {
                                    "sourceType": "tmdb",
                                    "movieId": 1,
                                    "title": "XXXX",
                                    "sourceId": 299534,
                                    "votes": 0,
                                    "voteCount": 0,
                                    "language": "english",
                                    "id": 1
                                }
                            ],
                            "secondaryYearSourceId": 0,
                            "sortTitle": "XXXX",
                            "sizeOnDisk": 0,
                            "status": "released",
                            "overview": "XXXX",
                            "inCinemas": "2019-04-23T22:00:00Z",
                            "physicalRelease": "2019-07-30T00:00:00Z",
                            "images": [
                                {
                                    "coverType": "poster",
                                    "url": "XXXX",
                                },
                                {
                                    "coverType": "fanart",
                                    "url": "XXXX",
                                }
                            ],
                            "website": "XXXX",
                            "downloaded": False,
                            "year": 2019,
                            "hasFile": False,
                            "youTubeTrailerId": "XXXX",
                            "studio": "XXXX",
                            "path": "XXXX",
                            "profileId": 1,
                            "pathState": "static",
                            "monitored": True,
                            "minimumAvailability": "announced",
                            "isAvailable": True,
                            "folderName": "XXXX",
                            "runtime": 181,
                            "lastInfoSync": "2020-04-22T19:36:38.549849Z",
                            "cleanTitle": "XXXX",
                            "imdbId": "XXXX",
                            "tmdbId": 1,
                            "titleSlug": "XXXX",
                            "genres": [],
                            "tags": [],
                            "added": "2020-04-22T19:36:38.289604Z",
                            "ratings": {
                                "votes": 12592,
                                "value": 8.3
                            },
                            "qualityProfileId": 6,
                            "id": 807
                        },
                        "quality": {
                            "quality": {
                                "id": 3,
                                "name": "WEBDL-1080p",
                                "source": "webdl",
                                "resolution": "r1080P",
                                "modifier": "none"
                            },
                            "customFormats": [],
                            "revision": {
                                "version": 1,
                                "real": 0
                            }
                        },
                        "size": 6691400202.0,
                        "title": "XXXX",
                        "sizeleft": 6392135680.0,
                        "timeleft": "01:25:09",
                        "estimatedCompletionTime": "2020-04-22T21:11:41.114137Z",
                        "status": "Downloading",
                        "trackedDownloadStatus": "Ok",
                        "statusMessages": [],
                        "downloadId": "EDDBD1F8E921FC00F97BF3FC4EAEA69693FBC198",
                        "protocol": "torrent",
                        "id": 1
                    },
                    {
                        "movie": {
                            "title": "XXXX",
                            "alternativeTitles": [
                                {
                                    "sourceType": "tmdb",
                                    "movieId": 1,
                                    "title": "XXXX",
                                    "sourceId": 1,
                                    "votes": 0,
                                    "voteCount": 0,
                                    "language": "english",
                                    "id": 1
                                }
                            ],
                            "secondaryYearSourceId": 0,
                            "sortTitle": "XXXX",
                            "sizeOnDisk": 0,
                            "status": "released",
                            "overview": "XXXX",
                            "inCinemas": "2018-04-24T22:00:00Z",
                            "physicalRelease": "2018-07-31T00:00:00Z",
                            "images": [
                                {
                                    "coverType": "poster",
                                    "url": "XXXX",
                                },
                                {
                                    "coverType": "fanart",
                                    "url": "XXXX",
                                }
                            ],
                            "website": "XXXX",
                            "downloaded": False,
                            "year": 2018,
                            "hasFile": False,
                            "youTubeTrailerId": "XXXX",
                            "studio": "XXXX",
                            "path": "XXX",
                            "profileId": 1,
                            "pathState": "static",
                            "monitored": True,
                            "minimumAvailability": "announced",
                            "isAvailable": True,
                            "folderName": "XXXX",
                            "runtime": 149,
                            "lastInfoSync": "2020-04-22T19:36:34.395332Z",
                            "cleanTitle": "avengersinfinitywar",
                            "imdbId": "XXXX",
                            "tmdbId": 1,
                            "titleSlug": "XXXX",
                            "genres": [],
                            "tags": [],
                            "added": "2020-04-22T19:36:34.090883Z",
                            "ratings": {
                                "votes": 17705,
                                "value": 8.3
                            },
                            "qualityProfileId": 6,
                            "id": 1
                        },
                        "quality": {
                            "quality": {
                                "id": 3,
                                "name": "WEBDL-1080p",
                                "source": "webdl",
                                "resolution": "r1080P",
                                "modifier": "none"
                            },
                            "customFormats": [],
                            "revision": {
                                "version": 1,
                                "real": 0
                            }
                        },
                        "size": 6874154725.0,
                        "title": "XXXX",
                        "sizeleft": 6786826240.0,
                        "timeleft": "07:53:47",
                        "estimatedCompletionTime": "2020-04-23T03:40:19.114172Z",
                        "status": "Downloading",
                        "trackedDownloadStatus": "Ok",
                        "statusMessages": [],
                        "downloadId": "BB82688F14D159A089AE7341CEFA10014AB0D56A",
                        "protocol": "torrent",
                        "id": 1428965449
                    }
                ])
            elif endpoint == 'api/movie':
                return FakeResponse(200, {}, [{
                    "title": "XXXX",
                    "alternativeTitles": [{
                        "sourceType": "tmdb", "movieId": 2,
                        "title": "XXXX",
                        "sourceId": 1,
                        "votes": 0,
                        "voteCount": 0,
                        "language": "english",
                        "id": 1
                    }],
                    "secondaryYearSourceId": 0,
                    "sortTitle": "XXXX",
                    "sizeOnDisk": 0,
                    "status": "announced",
                    "overview": "The plot is unknown at this time.",
                    "inCinemas": "2021-12-20T23:00:00Z",
                    "images": [{
                        "coverType": "poster",
                        "url": "XXXX",
                    }],
                    "website": "",
                    "downloaded": False,
                    "year": 2021,
                    "hasFile": False,
                    "youTubeTrailerId": "XXXX",
                    "path": "XXXX",
                    "profileId": 1,
                    "pathState": "static",
                    "monitored": True,
                    "minimumAvailability": "announced",
                    "isAvailable": True,
                    "folderName": "XXXX",
                    "runtime": 0,
                    "lastInfoSync": "2020-04-17T21:07:17.963456Z",
                    "cleanTitle": "XXXX",
                    "imdbId": "XXXX",
                    "tmdbId": 1,
                    "titleSlug": "XXXX",
                    "genres": [],
                    "tags": [],
                    "added": "2019-11-16T21:21:53.23528Z",
                    "ratings": {
                        "votes": 0,
                        "value": 0.0
                    },
                    "qualityProfileId": 1,
                    "id": 2
                }, {
                    "title": "XXXX",
                    "alternativeTitles": [],
                    "secondaryYearSourceId": 0,
                    "sortTitle": "clue",
                    "sizeOnDisk": 0,
                    "status": "announced",
                    "overview": "XXXX",
                    "images": [{
                        "coverType": "poster",
                        "url": "XXXX",
                    }],
                    "downloaded": False,
                    "year": 0,
                    "hasFile": False,
                    "studio": "XXXX",
                    "path": "XXXX",
                    "profileId": 1,
                    "pathState": "static",
                    "monitored": False,
                    "minimumAvailability": "announced",
                    "isAvailable": True,
                    "folderName": "XXXX",
                    "runtime": 0,
                    "lastInfoSync": "2020-04-17T21:07:01.150801Z",
                    "cleanTitle": "XXXX",
                    "imdbId": "XXXX",
                    "tmdbId": 1,
                    "titleSlug": "XXXX",
                    "genres": [],
                    "tags": [],
                    "added": "2019-11-16T21:22:25.391927Z",
                    "ratings": {
                        "votes": 0,
                        "value": 0.0
                    },
                    "qualityProfileId": 1,
                    "id": 4
                }])


def test_transmission():
    probe = TransmissionTestProbe('localhost:8080')
    probe.run()
    assert probe.measured()['activeTorrentCount'] == 1
    assert probe.measured()['pausedTorrentCount'] == 2
    assert probe.measured()['downloadSpeed'] == 1000
    assert probe.measured()['uploadSpeed'] == 500


def test_sonarr():
    probe = MonitorTestProbe('localhost:8080', MonitorProbe.App.sonarr, '')
    probe.run()
    measured = probe.measured()
    assert measured['calendar'] == 1
    assert measured['queue'] == 1
    assert measured['monitored'] == (1, 1)


def test_radarr():
    probe = MonitorTestProbe('localhost:8080', MonitorProbe.App.radarr, '')
    probe.run()
    measured = probe.measured()
    assert measured['calendar'] == 0
    assert measured['queue'] == 2
    assert measured['monitored'] == (1, 1)
