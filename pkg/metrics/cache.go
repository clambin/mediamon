package metrics

import (
	"github.com/clambin/go-common/http/roundtripper"
	"net/http"
)

var _ roundtripper.CacheMetrics = CustomizedCacheMetrics{}

type CustomizedCacheMetrics struct {
	roundtripper.CacheMetrics
	customize RequestCustomizer
}

func NewCustomizedCacheMetrics(namespace, subsystem, application string, f RequestCustomizer) roundtripper.CacheMetrics {
	return CustomizedCacheMetrics{
		CacheMetrics: roundtripper.NewCacheMetrics(namespace, subsystem, application),
		customize:    f,
	}
}

func (m CustomizedCacheMetrics) Measure(r *http.Request, found bool) {
	m.CacheMetrics.Measure(m.customize(r), found)
}
