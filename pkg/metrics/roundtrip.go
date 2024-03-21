package metrics

import (
	"github.com/clambin/go-common/http/roundtripper"
	"net/http"
	"time"
)

var _ roundtripper.RoundTripMetrics = CustomizedRoundTripMetrics{}

type CustomizedRoundTripMetrics struct {
	roundtripper.RoundTripMetrics
	customize RequestCustomizer
}

type RequestCustomizer func(r *http.Request) *http.Request

func NewCustomizedRoundTripMetrics(namespace, subsystem, application string, f RequestCustomizer) roundtripper.RoundTripMetrics {
	return CustomizedRoundTripMetrics{
		RoundTripMetrics: roundtripper.NewDefaultRoundTripMetrics(namespace, subsystem, application),
		customize:        f,
	}
}

func (m CustomizedRoundTripMetrics) Measure(req *http.Request, resp *http.Response, err error, duration time.Duration) {
	m.RoundTripMetrics.Measure(m.customize(req), resp, err, duration)
}
