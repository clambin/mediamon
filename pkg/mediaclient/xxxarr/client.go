package xxxarr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
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
	req.Header.Set(headers.AcceptEncoding, "identity")

	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("get %s: %w", target, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return response, &ErrHTTPFailed{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("read: %w", err)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		err = &ErrParseFailed{
			Err:  err,
			Body: body,
		}
	}
	return response, err
}
