package main

import (
	"github.com/clambin/mediamon/v2/internal/cmd/mediamon"
	"log/slog"
	"net/http"
	"os"

	_ "net/http/pprof"
)

var version = "change_me"

func main() {
	go func() {
		_ = http.ListenAndServe(":6060", nil)
	}()

	mediamon.RootCmd.Version = version
	if err := mediamon.RootCmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
