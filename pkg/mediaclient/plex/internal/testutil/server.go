package testutil

import "net/http"

func Handler(w http.ResponseWriter, req *http.Request) {
	if response, ok := plexResponses[req.URL.Path]; ok {
		_, _ = w.Write([]byte(response))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func WithToken(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Plex-Token") != token {
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		next(writer, request)
	}
}

var plexResponses = map[string]string{
	"/identity": `{ "MediaContainer": {
    	"size": 0,
    	"claimed": true,
    	"machineIdentifier": "SomeUUID",
    	"version": "SomeVersion"
  	}}`,
	"/status/sessions": `{ "MediaContainer": {
		"size": 2,
		"Metadata": [
			{ "User": { "title": "foo" },   "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "grandparentTitle": "series", "parentTitle": "season 1", "title": "pilot", "type": "episode"},
			{ "User": { "title": "bar" },   "Player": { "product": "Plex Web" }, "Session": { "location": "wan"}, "TranscodeSession": { "throttled": false, "videoDecision": "copy" }, "title": "movie 1" },
			{ "User": { "title": "snafu" }, "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "TranscodeSession": { "throttled": true, "speed": 3.1, "videoDecision": "transcode" }, "title": "movie 2" },
			{ "User": { "title": "snafu" }, "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "TranscodeSession": { "throttled": true, "speed": 4.1, "videoDecision": "transcode" }, "title": "movie 3" }
		]
	}}`,
	"/library/sections": `{ "MediaContainer": {
		"size": 2,
        "Directory": [
           { "Key": "1", "Type": "movie", "Title": "Movies" },
           { "Key": "2", "Type": "show", "Title": "Shows" }
        ]
    }}`,
	"/library/sections/1/all": `{ "MediaContainer" : {
        "Metadata": [
           { "guid": "1", "title": "foo" }
        ]
    }}`,
	"/library/sections/2/all": `{ "MediaContainer" : {
        "Metadata": [
           { "guid": "2", "title": "bar" }
        ]
    }}`,
}
