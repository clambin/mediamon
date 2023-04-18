package plex_test

import (
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    plex.Timestamp
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid",
			input:   "1655899131",
			want:    plex.Timestamp(time.Date(2022, time.June, 22, 11, 58, 51, 0, time.UTC)),
			wantErr: assert.NoError,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: assert.Error,
		},
		{
			name:    "invalid",
			input:   "abcd",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts plex.Timestamp
			tt.wantErr(t, (&ts).UnmarshalJSON([]byte(tt.input)))
			assert.Equal(t, tt.want, ts)
		})
	}
}