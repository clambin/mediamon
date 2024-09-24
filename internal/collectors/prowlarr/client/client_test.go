package client

import (
	"context"
	"errors"
	"github.com/clambin/mediaclients/prowlarr"
	"github.com/clambin/mediamon/v2/internal/collectors/prowlarr/client/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestProwlarr_GetIndexStats(t *testing.T) {
	ctx := context.Background()
	p := mocks.NewProwlarrClient(t)

	c, err := NewProwlarrClient("http://localhost", "token", http.DefaultClient)
	assert.NoError(t, err)
	c.Client = p

	p.EXPECT().
		GetApiV1IndexerstatsWithResponse(ctx, (*prowlarr.GetApiV1IndexerstatsParams)(nil)).
		Return(&prowlarr.GetApiV1IndexerstatsResponse{}, nil).
		Once()
	_, err = c.GetIndexStats(ctx)
	assert.NoError(t, err)

	p.EXPECT().
		GetApiV1IndexerstatsWithResponse(ctx, (*prowlarr.GetApiV1IndexerstatsParams)(nil)).
		Return(nil, errors.New("error")).
		Once()
	_, err = c.GetIndexStats(ctx)
	assert.Error(t, err)

}
