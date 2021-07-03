package mediaclient_test

import (
	"bytes"
	"context"
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestXXXArrClient_GetApplication(t *testing.T) {
	client := mediaclient.XXXArrClient{
		Client:      &http.Client{},
		URL:         "",
		APIKey:      "",
		Application: "sonarr",
	}

	assert.Equal(t, "sonarr", client.GetApplication(context.Background()))
}

func TestXXXArrClient_GetVersion(t *testing.T) {
	type fields struct {
		Client      *http.Client
		URL         string
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    "1.2.3.4444",
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    "1.2.3.4444",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mediaclient.XXXArrClient{
				Client:      tt.fields.Client,
				URL:         tt.fields.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetVersion(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetCalendar(t *testing.T) {
	type fields struct {
		Client      *http.Client
		URL         string
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mediaclient.XXXArrClient{
				Client:      tt.fields.Client,
				URL:         tt.fields.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetCalendar(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCalendar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCalendar() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetQueue(t *testing.T) {
	type fields struct {
		Client      *http.Client
		URL         string
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mediaclient.XXXArrClient{
				Client:      tt.fields.Client,
				URL:         tt.fields.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetQueue(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetMonitored(t *testing.T) {
	type fields struct {
		Client      *http.Client
		URL         string
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    2,
			want1:   1,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				Client:      httpstub.NewTestClient(xxxArrLoopback),
				URL:         "",
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    2,
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mediaclient.XXXArrClient{
				Client:      tt.fields.Client,
				URL:         tt.fields.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, got1, err := client.GetMonitored(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMonitored() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMonitored() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMonitored() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestXXXArrClient_GetMonitored_Panic(t *testing.T) {
	client := &mediaclient.XXXArrClient{
		Client:      httpstub.NewTestClient(xxxArrLoopback),
		URL:         "",
		APIKey:      "",
		Application: "invalid",
	}

	assert.Panics(t, func() { _, _, _ = client.GetMonitored(context.Background()) })
}

func TestXXXArrClient_ServerDown(t *testing.T) {
	client := &mediaclient.XXXArrClient{
		Client:      httpstub.NewTestClient(httpstub.Failing),
		URL:         "",
		APIKey:      "",
		Application: "sonarr",
	}

	_, err := client.GetVersion(context.Background())

	assert.NotNil(t, err)
	assert.Equal(t, "internal server error", err.Error())
}

// Server loopback function
func xxxArrLoopback(req *http.Request) *http.Response {
	if req.Header.Get("X-Api-Key") != "1234" {
		return &http.Response{
			StatusCode: 409,
			Status:     "No/invalid Application Key",
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}
	switch req.URL.Path {
	case "/api/v3/system/status":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(systemStatus)),
		}
	case "/api/v3/calendar":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(calendar)),
		}
	case "/api/v3/queue":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(queued)),
		}
	case "/api/v3/series":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(monitored)),
		}
	case "/api/v3/movie":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(monitored)),
		}
	default:
		return &http.Response{
			StatusCode: 404,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
		}
	}
}

// Responses

const (
	systemStatus = `{
  "version": "1.2.3.4444"
}`

	calendar = `[
  {
    "hasFile": false
  },
  {
    "hasFile": true
  }
]`
	queued = `{ "totalRecords": 3 }`

	monitored = `[ { "monitored": true }, { "monitored": false }, { "monitored": true } ]`
)
