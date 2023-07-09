package plex

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const authURL = "https://plex.tv/users/sign_in.xml"

// SetAuthToken sets the AuthToken
func (c *Client) SetAuthToken(s string) {
	c.plexAuth.lock.Lock()
	defer c.plexAuth.lock.Unlock()
	c.plexAuth.authToken = s
}

// GetAuthToken logs into plex.tv and returns the current authToken.
func (c *Client) GetAuthToken(ctx context.Context) (string, error) {
	err := c.plexAuth.authenticate(ctx)
	if err != nil {
		return "", err
	}

	c.plexAuth.lock.Lock()
	defer c.plexAuth.lock.Unlock()
	return c.plexAuth.authToken, nil
}

var _ http.RoundTripper = &authenticator{}

type authenticator struct {
	httpClient *http.Client
	username   string
	password   string
	authURL    string
	product    string
	version    string
	next       http.RoundTripper
	lock       sync.Mutex
	authToken  string
}

func (a *authenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := a.authenticate(request.Context()); err != nil {
		return nil, err
	}
	request.Header.Add("X-Plex-Token", a.authToken)

	return a.next.RoundTrip(request)
}

func (a *authenticator) authenticate(ctx context.Context) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.authToken != "" {
		return nil
	}

	req, _ := a.makeAuthRequest(ctx)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {
		a.authToken, err = getAuthResponse(resp.Body)
	} else {
		err = fmt.Errorf("plex auth: %s", resp.Status)
	}
	_ = resp.Body.Close()

	return err
}

func (a *authenticator) makeAuthRequest(ctx context.Context) (*http.Request, error) {
	v := make(url.Values)
	v.Set("user[login]", a.username)
	v.Set("user[password]", a.password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.authURL, bytes.NewBufferString(v.Encode()))
	if err == nil {
		req.Header.Add("X-Plex-Product", a.product)
		req.Header.Add("X-Plex-Version", a.version)
		req.Header.Add("X-Plex-Client-Identifier", a.product+"-v"+a.version)
	}
	return req, err
}

func getAuthResponse(body io.Reader) (string, error) {
	var authResponse struct {
		XMLName             xml.Name `xml:"user"`
		AuthenticationToken string   `xml:"authenticationToken,attr"`
	}
	err := xml.NewDecoder(body).Decode(&authResponse)
	return authResponse.AuthenticationToken, err
}
