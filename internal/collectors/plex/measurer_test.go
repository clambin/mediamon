package plex

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestMeasurer_Collect(t *testing.T) {
	type args struct {
		path     string
		err      error
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "metadata",
			args: args{
				path:     "/library/metadata/123/children",
				err:      nil,
				duration: 100 * time.Millisecond,
			},
			want: `
# HELP foo_bar_api_errors_total Number of failed HTTP calls
# TYPE foo_bar_api_errors_total counter
foo_bar_api_errors_total{application="test",method="",path="/library/metadata"} 0
# HELP foo_bar_api_latency latency of HTTP calls
# TYPE foo_bar_api_latency summary
foo_bar_api_latency_sum{application="test",method="",path="/library/metadata"} 0.1
foo_bar_api_latency_count{application="test",method="",path="/library/metadata"} 1
`,
		},
		{
			name: "sections",
			args: args{
				path:     "/library/sections/123/all",
				err:      nil,
				duration: 100 * time.Millisecond,
			},
			want: `
# HELP foo_bar_api_errors_total Number of failed HTTP calls
# TYPE foo_bar_api_errors_total counter
foo_bar_api_errors_total{application="test",method="",path="/library/sections"} 0
# HELP foo_bar_api_latency latency of HTTP calls
# TYPE foo_bar_api_latency summary
foo_bar_api_latency_sum{application="test",method="",path="/library/sections"} 0.1
foo_bar_api_latency_count{application="test",method="",path="/library/sections"} 1
`,
		},
		{
			name: "sessions",
			args: args{
				path:     "/status/sessions",
				err:      errors.New("failed"),
				duration: 100 * time.Millisecond,
			},
			want: `
# HELP foo_bar_api_errors_total Number of failed HTTP calls
# TYPE foo_bar_api_errors_total counter
foo_bar_api_errors_total{application="test",method="",path="/status/sessions"} 1
# HELP foo_bar_api_latency latency of HTTP calls
# TYPE foo_bar_api_latency summary
foo_bar_api_latency_sum{application="test",method="",path="/status/sessions"} 0.1
foo_bar_api_latency_count{application="test",method="",path="/status/sessions"} 1
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newMeasurer("foo", "bar", "test")
			r := prometheus.NewPedanticRegistry()
			require.NoError(t, r.Register(m))

			m.MeasureRequest(&http.Request{URL: &url.URL{Path: tt.args.path}}, nil, tt.args.err, tt.args.duration)
			assert.NoError(t, testutil.GatherAndCompare(r, strings.NewReader(tt.want)))
		})
	}
}
