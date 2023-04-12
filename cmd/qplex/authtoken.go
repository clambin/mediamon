package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func authToken(_ *cobra.Command, _ []string) {
	c := MakePlexClient()

	token, err := c.GetAuthToken(context.Background())
	if err != nil {
		slog.Error("failed to get authentication token", err)
		return
	}
	fmt.Printf("authToken: %s\n", token)
}
