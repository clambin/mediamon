package scraper

// Stats contains the statistics returned by a Scraper
type Stats struct {
	URL         string
	Version     string
	Calendar    []string
	Queued      []QueuedFile
	Monitored   int
	Unmonitored int
}

// QueuedFile contains information on a file that's queued for download
type QueuedFile struct {
	Name            string
	TotalBytes      float64
	DownloadedBytes float64
}

// Scraper provides a generic means of getting stats from Sonarr or Radarr
//go:generate mockery --name Scraper
type Scraper interface {
	Scrape() (Stats, error)
}
