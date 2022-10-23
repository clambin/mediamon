package xxxarr

import "time"

// RadarrSystemStatusResponse holds the response to Radarr's /api/v3/system/system endpoint
type RadarrSystemStatusResponse struct {
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
	IsNetCore              bool      `json:"isNetCore"`
	IsLinux                bool      `json:"isLinux"`
	IsOsx                  bool      `json:"isOsx"`
	IsWindows              bool      `json:"isWindows"`
	IsDocker               bool      `json:"isDocker"`
	Mode                   string    `json:"mode"`
	Branch                 string    `json:"branch"`
	Authentication         string    `json:"authentication"`
	DatabaseType           string    `json:"databaseType"`
	DatabaseVersion        string    `json:"databaseVersion"`
	MigrationVersion       int       `json:"migrationVersion"`
	URLBase                string    `json:"urlBase"`
	RuntimeVersion         string    `json:"runtimeVersion"`
	RuntimeName            string    `json:"runtimeName"`
	StartTime              time.Time `json:"startTime"`
	PackageVersion         string    `json:"packageVersion"`
	PackageAuthor          string    `json:"packageAuthor"`
	PackageUpdateMechanism string    `json:"packageUpdateMechanism"`
}

// RadarrHealthResponse holders the response of Radar's /api/v3/system/health endpoint
type RadarrHealthResponse struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiUrl string `json:"wikiUrl"`
}

// RadarrCalendarResponse holds the response of Radarr's /api/v3/calendar endpoint
type RadarrCalendarResponse struct {
	Title            string `json:"title"`
	OriginalTitle    string `json:"originalTitle"`
	OriginalLanguage struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"originalLanguage"`
	AlternateTitles []struct {
		SourceType      string `json:"sourceType"`
		MovieMetadataID int    `json:"movieMetadataId"`
		Title           string `json:"title"`
		SourceID        int    `json:"sourceId"`
		Votes           int    `json:"votes"`
		VoteCount       int    `json:"voteCount"`
		Language        struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		ID int `json:"id"`
	} `json:"alternateTitles"`
	SecondaryYearSourceID int       `json:"secondaryYearSourceId"`
	SortTitle             string    `json:"sortTitle"`
	SizeOnDisk            int       `json:"sizeOnDisk"`
	Status                string    `json:"status"`
	Overview              string    `json:"overview"`
	InCinemas             time.Time `json:"inCinemas"`
	Images                []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Website             string        `json:"website"`
	Year                int           `json:"year"`
	HasFile             bool          `json:"hasFile"`
	YouTubeTrailerID    string        `json:"youTubeTrailerId"`
	Studio              string        `json:"studio"`
	Path                string        `json:"path"`
	QualityProfileID    int           `json:"qualityProfileId"`
	Monitored           bool          `json:"monitored"`
	MinimumAvailability string        `json:"minimumAvailability"`
	IsAvailable         bool          `json:"isAvailable"`
	FolderName          string        `json:"folderName"`
	Runtime             int           `json:"runtime"`
	CleanTitle          string        `json:"cleanTitle"`
	ImdbID              string        `json:"imdbId"`
	TmdbID              int           `json:"tmdbId"`
	TitleSlug           string        `json:"titleSlug"`
	Certification       string        `json:"certification"`
	Genres              []string      `json:"genres"`
	Tags                []interface{} `json:"tags"`
	Added               time.Time     `json:"added"`
	Ratings             struct {
		Tmdb struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"tmdb"`
	} `json:"ratings"`
	Collection struct {
		Name   string        `json:"name"`
		TmdbID int           `json:"tmdbId"`
		Images []interface{} `json:"images"`
	} `json:"collection"`
	Popularity float64 `json:"popularity"`
	ID         int     `json:"id"`
}

// RadarrMovieResponse holds the response of Radarr's /api/v3/movie endpoint
type RadarrMovieResponse struct {
	Title            string `json:"title"`
	OriginalTitle    string `json:"originalTitle"`
	OriginalLanguage struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"originalLanguage"`
	AlternateTitles []struct {
		SourceType      string `json:"sourceType"`
		MovieMetadataID int    `json:"movieMetadataId"`
		Title           string `json:"title"`
		SourceID        int    `json:"sourceId"`
		Votes           int    `json:"votes"`
		VoteCount       int    `json:"voteCount"`
		Language        struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"language"`
		ID int `json:"id"`
	} `json:"alternateTitles"`
	SecondaryYearSourceID int    `json:"secondaryYearSourceId"`
	SortTitle             string `json:"sortTitle"`
	SizeOnDisk            int    `json:"sizeOnDisk"`
	Status                string `json:"status"`
	Overview              string `json:"overview"`
	Images                []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"images"`
	Website             string        `json:"website"`
	Year                int           `json:"year"`
	HasFile             bool          `json:"hasFile"`
	YouTubeTrailerID    string        `json:"youTubeTrailerId"`
	Studio              string        `json:"studio"`
	Path                string        `json:"path"`
	QualityProfileID    int           `json:"qualityProfileId"`
	Monitored           bool          `json:"monitored"`
	MinimumAvailability string        `json:"minimumAvailability"`
	IsAvailable         bool          `json:"isAvailable"`
	FolderName          string        `json:"folderName"`
	Runtime             int           `json:"runtime"`
	CleanTitle          string        `json:"cleanTitle"`
	ImdbID              string        `json:"imdbId,omitempty"`
	TmdbID              int           `json:"tmdbId"`
	TitleSlug           string        `json:"titleSlug"`
	Genres              []string      `json:"genres"`
	Tags                []interface{} `json:"tags"`
	Added               time.Time     `json:"added"`
	Ratings             struct {
		Tmdb struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"tmdb"`
		Imdb struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"imdb,omitempty"`
	} `json:"ratings"`
	Collection struct {
		Name   string        `json:"name"`
		TmdbID int           `json:"tmdbId"`
		Images []interface{} `json:"images"`
	} `json:"collection,omitempty"`
	Popularity     float64   `json:"popularity"`
	ID             int       `json:"id"`
	InCinemas      time.Time `json:"inCinemas,omitempty"`
	Certification  string    `json:"certification,omitempty"`
	DigitalRelease time.Time `json:"digitalRelease,omitempty"`
}

// RadarrQueueResponse contains the response to Radarr's /api/v3/queue endpoint
type RadarrQueueResponse struct {
	Page          int                         `json:"page"`
	PageSize      int                         `json:"pageSize"`
	SortKey       string                      `json:"sortKey"`
	SortDirection string                      `json:"sortDirection"`
	TotalRecords  int                         `json:"totalRecords"`
	Records       []RadarrQueueResponseRecord `json:"records"`
}

// RadarrQueueResponseRecord contains one record from a RadarrQueueResponse
type RadarrQueueResponseRecord struct {
	MovieID   int `json:"movieId"`
	Languages []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"languages"`
	Quality struct {
		Quality struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Source     string `json:"source"`
			Resolution int    `json:"resolution"`
			Modifier   string `json:"modifier"`
		} `json:"quality"`
		Revision struct {
			Version  int  `json:"version"`
			Real     int  `json:"real"`
			IsRepack bool `json:"isRepack"`
		} `json:"revision"`
	} `json:"quality"`
	CustomFormats           []interface{} `json:"customFormats"`
	Size                    int64         `json:"size"`
	Title                   string        `json:"title"`
	Sizeleft                int64         `json:"sizeleft"`
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
