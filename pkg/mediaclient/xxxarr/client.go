package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func call[T any](ctx context.Context, client *http.Client, target string) (T, error) {
	var response T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return response, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		err = &ErrInvalidJSON{
			Err:  err,
			Body: body,
		}
	}
	return response, err
}
