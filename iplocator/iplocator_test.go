package iplocator

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Locate(t *testing.T) {
	type want struct {
		err      assert.ErrorAssertionFunc
		location Location
	}
	tests := []struct {
		name    string
		address string
		want
	}{
		{
			name:    "valid",
			address: "8.8.8.8",
			want: want{
				err:      assert.NoError,
				location: Location{Status: "success", Lon: -77.5, Lat: 39.03},
			},
		},
		{
			name:    "invalid",
			address: "192.168.0.1",
			want:    want{err: assert.Error},
		},
		{
			name:    "unknown",
			address: "unknown",
			want:    want{err: assert.Error},
		},
	}

	s := server{responses: defaultResponses}
	ts := httptest.NewServer(&s)

	c := New(nil)
	c.url = ts.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := c.Locate(tt.address)
			tt.want.err(t, err)
			if err == nil {
				assert.Equal(t, tt.want.location, location)
			}
		})
	}

	assert.Equal(t, 3, s.calls)

	ts.Close()
	_, err := c.Locate("8.8.4.4")
	assert.Error(t, err)
}

var defaultResponses = map[string]Location{
	"/json/8.8.8.8": {
		Status: "success",
		Lon:    -77.5,
		Lat:    39.03,
	},
	"/json/192.168.0.1": {
		Status:  "fail",
		Message: "private range",
	},
}

type server struct {
	calls     int
	responses map[string]Location
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.calls++
	if resp, ok := s.responses[req.URL.Path]; ok {
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}
