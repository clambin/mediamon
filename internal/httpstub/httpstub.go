package httpstub

import (
	"net/http"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// Failing emulates a failing server. Always returns HTTP 500 error
func Failing(_ *http.Request) *http.Response {
	return &http.Response{
		StatusCode: 500,
		Status:     "internal server error",
	}
}
