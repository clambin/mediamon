package caller_test

import (
	"encoding/json"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/caller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCacher_Do(t *testing.T) {
	s := &server{}
	srv := httptest.NewServer(http.HandlerFunc(s.handle))
	c := caller.NewCacher(
		nil, "foo", caller.Options{},
		[]caller.CacheTableEntry{
			{Endpoint: "/foo"},
		},
		50*time.Millisecond, 0,
	)

	value, err := doCall2(c, srv.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, 1, value)

	value, err = doCall2(c, srv.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, 1, value)

	assert.Eventually(t, func() bool {
		value, err = doCall2(c, srv.URL+"/foo")
		return err == nil && value == 2
	}, 75*time.Millisecond, 10*time.Millisecond)

	value, err = doCall2(c, srv.URL+"/bar")
	require.NoError(t, err)
	assert.Equal(t, 3, value)

	value, err = doCall2(c, srv.URL+"/bar")
	require.NoError(t, err)
	assert.Equal(t, 4, value)

	srv.Close()

	value, err = doCall2(c, srv.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, 2, value)

	assert.Eventually(t, func() bool {
		_, err = doCall2(c, srv.URL+"/foo")
		return err != nil
	}, 75*time.Millisecond, 10*time.Millisecond)

}

func TestCacher_Do_MultipleEndpoints(t *testing.T) {
	s := &server{}
	srv := httptest.NewServer(http.HandlerFunc(s.handle))
	c := caller.NewCacher(
		nil, "foo", caller.Options{},
		[]caller.CacheTableEntry{
			{Endpoint: "/foo"},
			{Endpoint: "/bar", Expiry: 100 * time.Millisecond},
		},
		50*time.Millisecond, 0,
	)

	value, err := doCall2(c, srv.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, 1, value)

	value, err = doCall2(c, srv.URL+"/bar")
	require.NoError(t, err)
	assert.Equal(t, 2, value)

	value, err = doCall2(c, srv.URL+"/foo")
	require.NoError(t, err)
	assert.Equal(t, 1, value)

	value, err = doCall2(c, srv.URL+"/bar")
	require.NoError(t, err)
	assert.Equal(t, 2, value)

	assert.Eventually(t, func() bool {
		value, err = doCall2(c, srv.URL+"/foo")
		return err == nil && value == 3
	}, 75*time.Millisecond, 10*time.Millisecond)

	value, err = doCall2(c, srv.URL+"/bar")
	require.NoError(t, err)
	assert.Equal(t, 2, value)

	assert.Eventually(t, func() bool {
		value, err = doCall2(c, srv.URL+"/bar")
		return err == nil && value == 4
	}, 75*time.Millisecond, 10*time.Millisecond)

}

type server struct {
	counter int
}

type serverResponse struct {
	Counter int
}

func (s *server) handle(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/foo" && req.URL.Path != "/bar" {
		http.Error(w, "invalid endpoint: "+req.URL.Path, http.StatusNotFound)
		return
	}

	s.counter++
	err := json.NewEncoder(w).Encode(serverResponse{Counter: s.counter})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func doCall2(c caller.Caller, url string) (response int, err error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	var resp *http.Response
	if resp, err = c.Do(req); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("call failed: %s", resp.Status)
	}
	defer func() { _ = resp.Body.Close() }()

	var r serverResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	return r.Counter, err
}
