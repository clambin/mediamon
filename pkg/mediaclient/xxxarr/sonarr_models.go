package xxxarr

import "time"

// SonarrSystemStatusResponse contains the response of Sonarr's /api/v3/system/status endpoint
type SonarrSystemStatusResponse struct {
	AppName                string    `json:"appName"`
	Version                string    `json:"version"`
	BuildTime              time.Time `json:"buildTime"`
	IsDebug                bool      `json:"isDebug"`
	IsProduction           bool      `json:"isProduction"`
	IsAdmin                bool      `json:"isAdmin"`
	IsUserInteractive      bool      `json:"isUserInteractive"`
	StartupPath            string    `json:"startupPath"`
	AppData                string    `json:"appData"`
	OsName                 string    `json:"osName"`
	OsVersion              string    `json:"osVersion"`
	IsMonoRuntime          bool      `json:"isMonoRuntime"`
	IsMono                 bool      `json:"isMono"`
	IsLinux                bool      `json:"isLinux"`
	IsOsx                  bool      `json:"isOsx"`
	IsWindows              bool      `json:"isWindows"`
	Mode                   string    `json:"mode"`
	Branch                 string    `json:"branch"`
	Authentication         string    `json:"authentication"`
	SqliteVersion          string    `json:"sqliteVersion"`
	URLBase                string    `json:"urlBase"`
	RuntimeVersion         string    `json:"runtimeVersion"`
	RuntimeName            string    `json:"runtimeName"`
	StartTime              time.Time `json:"startTime"`
	PackageVersion         string    `json:"packageVersion"`
	PackageAuthor          string    `json:"packageAuthor"`
	PackageUpdateMechanism string    `json:"packageUpdateMechanism"`
}

// SonarrHealthResponse holders the response of Sonarr's /api/v3/system/health endpoint
type SonarrHealthResponse struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiUrl string `json:"wikiUrl"`
}

// SonarrCalendarResponse contains the response of Sonarr's /api/v3/calendar endpoint
type SonarrCalendarResponse struct {
	SeriesID                 int       `json:"seriesId"`
	TvdbID                   int       `json:"tvdbId"`
	EpisodeFileID            int       `json:"episodeFileId"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUtc               time.Time `json:"airDateUtc"`
	Overview                 string    `json:"overview"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	AbsoluteEpisodeNumber    int       `json:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
	ID                       int       `json:"id"`
}

// SonarrQueueResponse contains the response of Sonarr's /api/v3/queue endpoint
type SonarrQueueResponse struct {
	Page          int                         `json:"page"`
	PageSize      int                         `json:"pageSize"`
	SortKey       string                      `json:"sortKey"`
	SortDirection string                      `json:"sortDirection"`
	TotalRecords  int                         `json:"totalRecords"`
	Records       []SonarrQueueResponseRecord `json:"records"`
}

// SonarrQueueResponseRecord contains a Record from SonarrQueueResponse
type SonarrQueueResponseRecord struct {
	SeriesID  int `json:"seriesId"`
	EpisodeID int `json:"episodeId"`
	Language  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"language"`
	Quality struct {
		Quality struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Source     string `json:"source"`
			Resolution int    `json:"resolution"`
		} `json:"quality"`
		Revision struct {
			Version  int  `json:"version"`
			Real     int  `json:"real"`
			IsRepack bool `json:"isRepack"`
		} `json:"revision"`
	} `json:"quality"`
	Size                    float64       `json:"size"`
	Title                   string        `json:"title"`
	Sizeleft                float64       `json:"sizeleft"`
	Timeleft                string        `json:"timeleft"`
	EstimatedCompletionTime time.Time     `json:"estimatedCompletionTime"`
	Status                  string        `json:"status"`
	TrackedDownloadStatus   string        `json:"trackedDownloadStatus"`
	TrackedDownloadState    string        `json:"trackedDownloadState"`
	StatusMessages          []interface{} `json:"statusMessages"`
	DownloadID              string        `json:"downloadId"`
	Protocol                string        `json:"protocol"`
	DownloadClient          string        `json:"downloadClient"`
	Indexer                 string        `json:"indexer"`
	OutputPath              string        `json:"outputPath"`
	ID                      int           `json:"id"`
}

// SonarrSeriesResponse contains the response of Sonarr's /api/v3/series endpoint
type SonarrSeriesResponse struct {
	Title           string        `json:"title"`
	AlternateTitles []interface{} `json:"alternateTitles"`
	SortTitle       string        `json:"sortTitle"`
	Status          string        `json:"status"`
	Ended           bool          `json:"ended"`
	Overview        string        `json:"overview"`
	PreviousAiring  time.Time     `json:"previousAiring"`
	Network         string        `json:"network"`
	AirTime         string        `json:"airTime"`
	Images          []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
		Statistics   struct {
			EpisodeFileCount  int       `json:"episodeFileCount"`
			EpisodeCount      int       `json:"episodeCount"`
			TotalEpisodeCount int       `json:"totalEpisodeCount"`
			SizeOnDisk        int64     `json:"sizeOnDisk"`
			ReleaseGroups     []string  `json:"releaseGroups"`
			PercentOfEpisodes float64   `json:"percentOfEpisodes"`
			PreviousAiring    time.Time `json:"previousAiring,omitempty"`
		} `json:"statistics"`
	} `json:"seasons"`
	Year              int           `json:"year"`
	Path              string        `json:"path"`
	QualityProfileID  int           `json:"qualityProfileId"`
	LanguageProfileID int           `json:"languageProfileId"`
	SeasonFolder      bool          `json:"seasonFolder"`
	Monitored         bool          `json:"monitored"`
	UseSceneNumbering bool          `json:"useSceneNumbering"`
	Runtime           int           `json:"runtime"`
	TvdbID            int           `json:"tvdbId"`
	TvRageID          int           `json:"tvRageId"`
	TvMazeID          int           `json:"tvMazeId"`
	FirstAired        time.Time     `json:"firstAired"`
	SeriesType        string        `json:"seriesType"`
	CleanTitle        string        `json:"cleanTitle"`
	ImdbID            string        `json:"imdbId"`
	TitleSlug         string        `json:"titleSlug"`
	RootFolderPath    string        `json:"rootFolderPath"`
	Certification     string        `json:"certification"`
	Genres            []string      `json:"genres"`
	Tags              []interface{} `json:"tags"`
	Added             time.Time     `json:"added"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	Statistics struct {
		SeasonCount       int      `json:"seasonCount"`
		EpisodeFileCount  int      `json:"episodeFileCount"`
		EpisodeCount      int      `json:"episodeCount"`
		TotalEpisodeCount int      `json:"totalEpisodeCount"`
		SizeOnDisk        int64    `json:"sizeOnDisk"`
		ReleaseGroups     []string `json:"releaseGroups"`
		PercentOfEpisodes float64  `json:"percentOfEpisodes"`
	} `json:"statistics"`
	ID int `json:"id"`
}

// SonarrEpisodeResponse contains the response from Sonarr's /api/v3/episode endpoint
type SonarrEpisodeResponse struct {
	SeriesID      int       `json:"seriesId"`
	TvdbID        int       `json:"tvdbId"`
	EpisodeFileID int       `json:"episodeFileId"`
	SeasonNumber  int       `json:"seasonNumber"`
	EpisodeNumber int       `json:"episodeNumber"`
	Title         string    `json:"title"`
	AirDate       string    `json:"airDate"`
	AirDateUtc    time.Time `json:"airDateUtc"`
	Overview      string    `json:"overview"`
	EpisodeFile   struct {
		SeriesID     int       `json:"seriesId"`
		SeasonNumber int       `json:"seasonNumber"`
		RelativePath string    `json:"relativePath"`
		Path         string    `json:"path"`
		Size         int       `json:"size"`
		DateAdded    time.Time `json:"dateAdded"`
		SceneName    string    `json:"sceneName"`
		ReleaseGroup string    `json:"releaseGroup"`
		Language     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		Quality struct {
			Quality struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Source     string `json:"source"`
				Resolution int    `json:"resolution"`
			} `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		MediaInfo struct {
			AudioBitrate     int     `json:"audioBitrate"`
			AudioChannels    float64 `json:"audioChannels"`
			AudioCodec       string  `json:"audioCodec"`
			AudioLanguages   string  `json:"audioLanguages"`
			AudioStreamCount int     `json:"audioStreamCount"`
			VideoBitDepth    int     `json:"videoBitDepth"`
			VideoBitrate     int     `json:"videoBitrate"`
			VideoCodec       string  `json:"videoCodec"`
			VideoFps         float64 `json:"videoFps"`
			Resolution       string  `json:"resolution"`
			RunTime          string  `json:"runTime"`
			ScanType         string  `json:"scanType"`
			Subtitles        string  `json:"subtitles"`
		} `json:"mediaInfo"`
		QualityCutoffNotMet  bool `json:"qualityCutoffNotMet"`
		LanguageCutoffNotMet bool `json:"languageCutoffNotMet"`
		ID                   int  `json:"id"`
	} `json:"episodeFile"`
	HasFile                  bool                        `json:"hasFile"`
	Monitored                bool                        `json:"monitored"`
	AbsoluteEpisodeNumber    int                         `json:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool                        `json:"unverifiedSceneNumbering"`
	Series                   SonarrEpisodeResponseSeries `json:"series"`
	Images                   []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	ID int `json:"id"`
}

// SonarrEpisodeResponseSeries contains the Series in a SonarrEpisodeResponse
type SonarrEpisodeResponseSeries struct {
	Title     string `json:"title"`
	SortTitle string `json:"sortTitle"`
	Status    string `json:"status"`
	Ended     bool   `json:"ended"`
	Overview  string `json:"overview"`
	Network   string `json:"network"`
	AirTime   string `json:"airTime"`
	Images    []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
	} `json:"seasons"`
	Year              int           `json:"year"`
	Path              string        `json:"path"`
	QualityProfileID  int           `json:"qualityProfileId"`
	LanguageProfileID int           `json:"languageProfileId"`
	SeasonFolder      bool          `json:"seasonFolder"`
	Monitored         bool          `json:"monitored"`
	UseSceneNumbering bool          `json:"useSceneNumbering"`
	Runtime           int           `json:"runtime"`
	TvdbID            int           `json:"tvdbId"`
	TvRageID          int           `json:"tvRageId"`
	TvMazeID          int           `json:"tvMazeId"`
	FirstAired        time.Time     `json:"firstAired"`
	SeriesType        string        `json:"seriesType"`
	CleanTitle        string        `json:"cleanTitle"`
	ImdbID            string        `json:"imdbId"`
	TitleSlug         string        `json:"titleSlug"`
	Certification     string        `json:"certification"`
	Genres            []string      `json:"genres"`
	Tags              []interface{} `json:"tags"`
	Added             time.Time     `json:"added"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	ID int `json:"id"`
}
