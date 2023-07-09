package plex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       Timestamp
		wantErr    assert.ErrorAssertionFunc
		wantString string
	}{
		{
			name:       "valid",
			input:      "1655899131",
			want:       Timestamp(time.Date(2022, time.June, 22, 11, 58, 51, 0, time.UTC)),
			wantErr:    assert.NoError,
			wantString: "2022-06-22 11:58:51 +0000 UTC",
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
			var ts Timestamp
			tt.wantErr(t, (&ts).UnmarshalJSON([]byte(tt.input)))
			assert.Equal(t, tt.want, ts)
			if tt.wantString != "" {
				assert.Equal(t, tt.wantString, ts.String())
			}
		})
	}
}
