package xxxarr

import "time"

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
