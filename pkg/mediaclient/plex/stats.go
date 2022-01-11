package plex

type identityStats struct {
	MediaContainer struct {
		Version string
	}
}

type sessionStats struct {
	MediaContainer struct {
		Metadata []struct {
			GrandparentTitle string
			Media            []struct {
				Part []struct {
					Stream []struct {
						Decision string
						Location string
					}
				}
			}
			User struct {
				Title string
			}
			Player struct {
				Local bool
			}
			TranscodeSession struct {
				Throttled     bool
				Speed         float64
				VideoDecision string
			}
		}
	}
}
