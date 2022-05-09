package updater

// Stats contains the statistics returned by a StatsGetter
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

// StatsGetter provides a generic means of getting stats from Sonarr or Radarr
type StatsGetter interface {
	GetStats() (Stats, error)
}
