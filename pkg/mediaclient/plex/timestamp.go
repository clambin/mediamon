package plex

import (
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(buf []byte) error {
	epoch, err := strconv.Atoi(string(buf))
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	*t = Timestamp(time.Unix(int64(epoch), 0).UTC())
	return nil
}

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}
