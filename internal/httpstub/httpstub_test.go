package httpstub_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mediamon/internal/httpstub"
	"net/http"
	"testing"
)

func TestNewTestClient(t *testing.T) {
	client := httpstub.NewTestClient(loopback)

	const message = "Hello world"

	reqBody := bytes.NewBufferString(message)
	req, _ := http.NewRequest("GET", "", reqBody)

	resp, err := client.Do(req)

	assert.Nil(t, err)

	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)

		assert.Nil(t, err)
		assert.Equal(t, message, string(respBody))
	}
}

func TestFailing(t *testing.T) {
	client := httpstub.NewTestClient(httpstub.Failing)

	const message = "Hello world"

	reqBody := bytes.NewBufferString(message)
	req, _ := http.NewRequest("GET", "", reqBody)

	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, "internal server error", resp.Status)
}

func loopback(req *http.Request) *http.Response {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err == nil {
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
		}
	}

	return &http.Response{
		StatusCode: 500,
		Status:     err.Error(),
		Header:     nil,
		Body:       nil,
	}
}
