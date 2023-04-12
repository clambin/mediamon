package plex

import "context"

// Identity contains the response of Plex's /identity API
type Identity struct {
	Size              int    `json:"size"`
	Claimed           bool   `json:"claimed"`
	MachineIdentifier string `json:"machineIdentifier"`
	Version           string `json:"version"`
}

// GetIdentity calls Plex' /identity endpoint. Mainly useful to get the server's version.
func (c *Client) GetIdentity(ctx context.Context) (identity Identity, err error) {
	err = c.call(ctx, "/identity", &identity)
	return
}
