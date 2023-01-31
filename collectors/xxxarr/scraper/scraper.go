package scraper

// Stats contains the statistics returned by a Scraper
type Stats struct {
	URL         string
	Version     string
	Health      map[string]int
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
