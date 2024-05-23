package iplocator

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Locate(t *testing.T) {
	s := server{
		responses: map[string]ipAPIResponse{
			"/json/8.8.8.8": {
				Status: "success",
				Lon:    -77.5,
				Lat:    39.03,
			},
			"/json/192.168.0.1": {
				Status:  "fail",
				Message: "private range",
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(s.handle))

	c := New(nil)
	c.url = ts.URL

	testCases := []struct {
		name    string
		address string
		wantErr assert.ErrorAssertionFunc
		wantLon float64
		wantLat float64
	}{
		{
			name:    "valid",
			address: "8.8.8.8",
			wantErr: assert.NoError,
			wantLon: -77.5,
			wantLat: 39.03,
		},
		{
			name:    "invalid",
			address: "192.168.0.1",
			wantErr: assert.Error,
		},
		{
			name:    "unknown",
			address: "unknown",
			wantErr: assert.Error,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			lon, lat, err := c.Locate(tt.address)
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.wantLon, lon)
				assert.Equal(t, tt.wantLat, lat)
			}
		})
	}

	assert.Equal(t, 3, s.calls)

	ts.Close()
	_, _, err := c.Locate("8.8.4.4")
	assert.Error(t, err)
}

type server struct {
	calls     int
	responses map[string]ipAPIResponse
}

func (s *server) handle(w http.ResponseWriter, req *http.Request) {
	s.calls++
	resp, ok := s.responses[req.URL.Path]
	if ok == false {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}
