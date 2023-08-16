package iplocator

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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
	c := New()
	c.url = ts.URL

	c.logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	lon, lat, err := c.Locate("8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, -77.5, lon)
	assert.Equal(t, 39.03, lat)

	lon, lat, err = c.Locate("8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, -77.5, lon)
	assert.Equal(t, 39.03, lat)

	_, _, err = c.Locate("192.168.0.1")
	assert.Error(t, err)

	assert.Equal(t, 2, s.calls)

	_, _, err = c.Locate("invalid")
	assert.Error(t, err)

	ts.Close()
	_, _, err = c.Locate("8.8.4.4")
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
