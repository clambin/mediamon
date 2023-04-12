package qplex

type viewCount map[guid]ViewCountEntry

func (vc viewCount) merge(s viewCount) {
	for id, entry := range s {
		if existing, ok := vc[id]; ok {
			entry.Views += existing.Views
		}
		vc[id] = entry
	}
}

func (vc viewCount) flatten() []ViewCountEntry {
	viewCounts := make([]ViewCountEntry, 0, len(vc))
	for _, entry := range vc {
		viewCounts = append(viewCounts, entry)
	}
	return viewCounts
}

type guid string
type ViewCountEntry struct {
	Library string
	Title   string
	Views   int
}
