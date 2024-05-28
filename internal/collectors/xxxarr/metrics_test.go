package xxxarr

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func Test_chopPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "unaltered",
			path: "/healthz",
			want: "/healthz",
		},
		{
			name: "unmatched",
			path: "/healthz/foo",
			want: "/healthz/foo",
		},
		{
			name: "altered",
			path: "/healthz/123",
			want: "/healthz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out := choppedPath(&http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: tt.path},
			})
			assert.Equal(t, tt.want, out)
		})
	}
}
