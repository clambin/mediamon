package prowlarr

import "github.com/prometheus/client_golang/prometheus"

func newMetrics(url string) map[string]*prometheus.Desc {
	constLabels := prometheus.Labels{"application": "prowlarr", "url": url}
	return map[string]*prometheus.Desc{
		"indexerResponseTime": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_response_time"),
			"Average response time in seconds",
			[]string{"indexer"},
			constLabels,
		),
		"indexerQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_query_total"),
			"Total number of queries to this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_grab_total"),
			"Total number of grabs from this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerFailedQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_failed_query_total"),
			"Total number of failed queries to this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"indexerFailedGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "indexer_failed_grab_total"),
			"Total number of failed grabs from this indexer",
			[]string{"indexer"},
			constLabels,
		),
		"userAgentQueryTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "user_agent_query_total"),
			"Total number of queries by user agent",
			[]string{"user_agent"},
			constLabels,
		),
		"userAgentGrabTotal": prometheus.NewDesc(
			prometheus.BuildFQName("mediamon", "prowlarr", "user_agent_grab_total"),
			"Total number of grabs by user agent",
			[]string{"user_agent"},
			constLabels,
		),
	}
}
