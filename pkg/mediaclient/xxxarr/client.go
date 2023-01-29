package xxxarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	var body bytes.Buffer
	r := io.TeeReader(resp.Body, &body)

	if err = json.NewDecoder(r).Decode(&response); err != nil {
		err = &ErrParseFailed{
			Err:  err,
			Body: body.Bytes(),
		}
	}
	return response, err
}
