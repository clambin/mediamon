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

const AuthURL = "https://plex.tv/users/sign_in.xml"

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

var _ http.RoundTripper = &Auth{}

type Auth struct {
	HTTPClient *http.Client
	Username   string
	Password   string
	AuthURL    string
	Product    string
	Version    string
	Next       http.RoundTripper
	authToken  string
	lock       sync.Mutex
}

func (a *Auth) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := a.authenticate(request.Context()); err != nil {
		return nil, err
	}
	request.Header.Add("X-Plex-Token", a.authToken)

	return a.Next.RoundTrip(request)
}

func (a *Auth) authenticate(ctx context.Context) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.authToken != "" {
		return nil
	}

	req, _ := a.makeAuthRequest(ctx)
	resp, err := a.HTTPClient.Do(req)
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

func (a *Auth) makeAuthRequest(ctx context.Context) (*http.Request, error) {
	v := make(url.Values)
	v.Set("user[login]", a.Username)
	v.Set("user[password]", a.Password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.AuthURL, bytes.NewBufferString(v.Encode()))
	if err == nil {
		req.Header.Add("X-Plex-Product", a.Product)
		req.Header.Add("X-Plex-Version", a.Version)
		req.Header.Add("X-Plex-Client-Identifier", a.Product+"-v"+a.Version)
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
