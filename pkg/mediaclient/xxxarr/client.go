package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func call[T any](ctx context.Context, client *http.Client, target, key string) (T, error) {
	var response T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return response, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("get %s: %w", target, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("call failed: " + resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		err = fmt.Errorf("decode %s: %w", target, err)
	}
	return response, err
}
