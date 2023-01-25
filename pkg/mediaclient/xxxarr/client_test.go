package xxxarr_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type TestServer struct {
	apiKey    string
	responses Responses
	server    *httptest.Server
}

func NewTestServer(responses Responses, apiKey string) *TestServer {
	s := &TestServer{
		apiKey:    apiKey,
		responses: responses,
	}
	s.server = httptest.NewServer(http.HandlerFunc(s.Handler))
	return s
}

type Responses map[string]any

func (ts TestServer) Handler(w http.ResponseWriter, req *http.Request) {
	// check auth
	if req.Header.Get("X-Api-Key") != ts.apiKey {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	endpoint := req.URL.Path
	if req.URL.RawQuery != "" {
		endpoint += "?" + req.URL.RawQuery
	}

	response, ok := ts.responses[endpoint]
	if !ok {
		http.Error(w, "endpoint not implemented", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode output: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
