package xxxarr

import "net/http"

var _ http.RoundTripper = &authentication{}

type authentication struct {
	apiKey string
	next   http.RoundTripper
}

func (a *authentication) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("X-Api-Key", a.apiKey)
	return a.next.RoundTrip(request)
}
