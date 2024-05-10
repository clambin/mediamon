package main

import (
	"github.com/clambin/mediamon/v2/internal/cmd/mediamon"
	"log/slog"
	"os"
)

var version = "change_me"

func main() {
	mediamon.RootCmd.Version = version
	if err := mediamon.RootCmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
