package clients

type QueuedItem struct {
	Name            string
	TotalBytes      int64
	DownloadedBytes int64
}

type Library struct {
	Monitored   int
	Unmonitored int
}
