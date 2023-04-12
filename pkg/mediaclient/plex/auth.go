package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clambin/mediamon/version"
	"io"
	"net/http"
	"net/url"
)

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
	req.Header.Add("X-Plex-Version", version.BuildVersion)
	req.Header.Add("X-Plex-Client-Identifier", "github.com/clambin/mediamon-v"+version.BuildVersion)

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

	var token string
	err := xml.Unmarshal(body, &authResponse)
	if err == nil {
		token = authResponse.AuthenticationToken
	}

	return token, err
}

func (c *Client) GetAccessTokens(ctx context.Context) ([]AccessToken, error) {
	var tokens []AccessToken
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://plex.tv/api/v2/server/access_tokens?auth_token="+c.AuthToken, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return tokens, err
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&tokens)
	return tokens, err
}
