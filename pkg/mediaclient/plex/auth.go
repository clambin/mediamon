package plex

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SetAuthToken sets the AuthToken
func (c *Client) SetAuthToken(s string) {
	c.AuthToken = s
}

// GetAuthToken logs into plex.tv and returns the current authToken.
func (c *Client) GetAuthToken(ctx context.Context) (string, error) {
	err := c.authenticate(ctx)
	if err != nil {
		return "", err
	}
	return c.AuthToken, nil
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (c *Client) authenticate(ctx context.Context) error {
	if c.AuthToken != "" {
		return nil
	}

	authURL := "https://plex.tv/users/sign_in.xml"
	if c.AuthURL != "" {
		authURL = c.AuthURL
	}
	authBody := c.makeAuthBody()

	product := c.Product
	if product == "" {
		product = "github.com/clambin/mediamon"
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, authURL, bytes.NewBufferString(authBody))
	req.Header.Add("X-Plex-Product", product)
	req.Header.Add("X-Plex-Version", c.Version)
	req.Header.Add("X-Plex-Client-Identifier", "github.com/clambin/mediamon-v"+c.Version)

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("plex auth failed: %s", resp.Status)
	}

	c.AuthToken, err = getAuthResponse(body)
	return err
}

func (c *Client) makeAuthBody() string {
	v := make(url.Values)
	v.Set("user[login]", c.UserName)
	v.Set("user[password]", c.Password)

	return v.Encode()
}

func getAuthResponse(body []byte) (string, error) {
	var authResponse struct {
		XMLName             xml.Name `xml:"user"`
		AuthenticationToken string   `xml:"authenticationToken,attr"`
	}

	err := xml.Unmarshal(body, &authResponse)
	return authResponse.AuthenticationToken, err
}
