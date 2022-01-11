package main

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"net/http"
	"os"
)

func main() {
	client := plex.Client{
		Client:  http.DefaultClient,
		URL:     "http://plex.192.168.0.11.nip.io",
		Options: plex.Options{},
		//AuthURL:  "",
		UserName: os.Getenv("PLEX_USER"),
		Password: os.Getenv("PLEX_PASSWORD"),
	}

	_, err := client.GetVersion(context.Background())
	if err != nil {
		panic(err)
	}

	sessions, err := client.GetSessions(context.Background())
	if err != nil {
		panic(err)
	}

	for _, session := range sessions {
		fmt.Printf("user: %s. transcode: %v (speed: %.1f, throttled: %v). media: %s. local: %v\n",
			session.User,
			session.Transcode, session.Speed, session.Throttled,
			session.Title,
			session.Local,
		)
	}
}
