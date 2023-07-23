package plex

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	"sort"
	"strings"
)

// Sessions contains the response of Plex's /status/sessions API
type Sessions struct {
	Size     int       `json:"size"`
	Metadata []Session `json:"Metadata"`
}

// GetSessions retrieves session information from the server.
func (c *Client) GetSessions(ctx context.Context) (sessions Sessions, err error) {
	err = c.call(ctx, "/status/sessions", &sessions)
	return
}

// Session contains one record in a Sessions
type Session struct {
	AddedAt               int            `json:"addedAt"`
	Art                   string         `json:"art"`
	AudienceRating        float64        `json:"audienceRating"`
	AudienceRatingImage   string         `json:"audienceRatingImage"`
	ContentRating         string         `json:"contentRating"`
	Duration              int            `json:"duration"`
	GrandparentArt        string         `json:"grandparentArt"`
	GrandparentGUID       string         `json:"grandparentGuid"`
	GrandparentKey        string         `json:"grandparentKey"`
	GrandparentRatingKey  string         `json:"grandparentRatingKey"`
	GrandparentTheme      string         `json:"grandparentTheme"`
	GrandparentThumb      string         `json:"grandparentThumb"`
	GrandparentTitle      string         `json:"grandparentTitle"`
	GUID                  string         `json:"guid"`
	Index                 int            `json:"index"`
	Key                   string         `json:"key"`
	LastViewedAt          Timestamp      `json:"lastViewedAt"`
	LibrarySectionID      string         `json:"librarySectionID"`
	LibrarySectionKey     string         `json:"librarySectionKey"`
	LibrarySectionTitle   string         `json:"librarySectionTitle"`
	OriginallyAvailableAt string         `json:"originallyAvailableAt"`
	ParentGUID            string         `json:"parentGuid"`
	ParentIndex           int            `json:"parentIndex"`
	ParentKey             string         `json:"parentKey"`
	ParentRatingKey       string         `json:"parentRatingKey"`
	ParentThumb           string         `json:"parentThumb"`
	ParentTitle           string         `json:"parentTitle"`
	Rating                float64        `json:"rating"`
	RatingKey             string         `json:"ratingKey"`
	SessionKey            string         `json:"sessionKey"`
	Summary               string         `json:"summary"`
	Thumb                 string         `json:"thumb"`
	Title                 string         `json:"title"`
	Type                  string         `json:"type"`
	UpdatedAt             Timestamp      `json:"updatedAt"`
	ViewOffset            int            `json:"viewOffset"`
	Media                 []SessionMedia `json:"Media"`
	Director              []struct {
		Filter string `json:"filter"`
		ID     string `json:"id"`
		Tag    string `json:"tag"`
	} `json:"Director"`
	Writer []struct {
		Filter string `json:"filter"`
		ID     string `json:"id"`
		Tag    string `json:"tag"`
	} `json:"Writer"`
	Rating2 []struct {
		Image string `json:"image"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"Rating"`
	Role []struct {
		Filter string `json:"filter"`
		ID     string `json:"id"`
		Role   string `json:"role"`
		Tag    string `json:"tag"`
		Thumb  string `json:"thumb,omitempty"`
	} `json:"Role"`
	User             SessionUser       `json:"User"`
	Player           SessionPlayer     `json:"Player"`
	Session          SessionStats      `json:"Session"`
	TranscodeSession SessionTranscoder `json:"TranscodeSession"`
}

// SessionMedia contains one record in a Session's Media list
type SessionMedia struct {
	AudioProfile          string             `json:"audioProfile"`
	ID                    string             `json:"id"`
	VideoProfile          string             `json:"videoProfile"`
	AudioChannels         int                `json:"audioChannels"`
	AudioCodec            string             `json:"audioCodec"`
	Bitrate               int                `json:"bitrate"`
	Container             string             `json:"container"`
	Duration              int                `json:"duration"`
	Height                int                `json:"height"`
	OptimizedForStreaming bool               `json:"optimizedForStreaming"`
	Protocol              string             `json:"protocol"`
	VideoCodec            string             `json:"videoCodec"`
	VideoFrameRate        string             `json:"videoFrameRate"`
	VideoResolution       string             `json:"videoResolution"`
	Width                 int                `json:"width"`
	Selected              bool               `json:"selected"`
	Part                  []MediaSessionPart `json:"Part"`
}

// MediaSessionPart contains one record in a MediaSession's Part list
type MediaSessionPart struct {
	AudioProfile          string                   `json:"audioProfile"`
	ID                    string                   `json:"id"`
	VideoProfile          string                   `json:"videoProfile"`
	Bitrate               int                      `json:"bitrate"`
	Container             string                   `json:"container"`
	Duration              int                      `json:"duration"`
	Height                int                      `json:"height"`
	OptimizedForStreaming bool                     `json:"optimizedForStreaming"`
	Protocol              string                   `json:"protocol"`
	Width                 int                      `json:"width"`
	Decision              string                   `json:"decision"`
	Selected              bool                     `json:"selected"`
	Stream                []MediaSessionPartStream `json:"Stream"`
}

// MediaSessionPartStream contains one stream (video, audio, subtitles) in a MediaSession's Part list
type MediaSessionPartStream struct {
	Bitrate              int     `json:"bitrate,omitempty"`
	Codec                string  `json:"codec"`
	Default              bool    `json:"default"`
	DisplayTitle         string  `json:"displayTitle"`
	ExtendedDisplayTitle string  `json:"extendedDisplayTitle"`
	FrameRate            float64 `json:"frameRate,omitempty"`
	Height               int     `json:"height,omitempty"`
	ID                   string  `json:"id"`
	Language             string  `json:"language"`
	LanguageCode         string  `json:"languageCode"`
	LanguageTag          string  `json:"languageTag"`
	StreamType           int     `json:"streamType"`
	Width                int     `json:"width,omitempty"`
	Decision             string  `json:"decision"`
	Location             string  `json:"location"`
	AudioChannelLayout   string  `json:"audioChannelLayout,omitempty"`
	BitrateMode          string  `json:"bitrateMode,omitempty"`
	Channels             int     `json:"channels,omitempty"`
	Profile              string  `json:"profile,omitempty"`
	SamplingRate         int     `json:"samplingRate,omitempty"`
	Selected             bool    `json:"selected,omitempty"`
	Title                string  `json:"title,omitempty"`
	Container            string  `json:"container,omitempty"`
	Format               string  `json:"format,omitempty"`
}

// SessionUser contains the user details inside a Session
type SessionUser struct {
	ID    string `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}

// SessionPlayer contains the player details inside a Session
type SessionPlayer struct {
	Address             string `json:"address"`
	Device              string `json:"device"`
	MachineIdentifier   string `json:"machineIdentifier"`
	Model               string `json:"model"`
	Platform            string `json:"platform"`
	PlatformVersion     string `json:"platformVersion"`
	Product             string `json:"product"`
	Profile             string `json:"profile"`
	RemotePublicAddress string `json:"remotePublicAddress"`
	State               string `json:"state"`
	Title               string `json:"title"`
	Version             string `json:"version"`
	Local               bool   `json:"local"`
	Relayed             bool   `json:"relayed"`
	Secure              bool   `json:"secure"`
	UserID              int    `json:"userID"`
}

// SessionStats contains the session details inside a Session
type SessionStats struct {
	ID        string `json:"id"`
	Bandwidth int    `json:"bandwidth"`
	Location  string `json:"location"`
}

// SessionTranscoder contains the transcoder details inside a Session.
// If the session doesn't transcode any media streams, all fields will be blank.
type SessionTranscoder struct {
	Key                     string  `json:"key"`
	Throttled               bool    `json:"throttled"`
	Complete                bool    `json:"complete"`
	Progress                float64 `json:"progress"`
	Size                    int     `json:"size"`
	Speed                   float64 `json:"speed"`
	Error                   bool    `json:"error"`
	Duration                int     `json:"duration"`
	Context                 string  `json:"context"`
	SourceVideoCodec        string  `json:"sourceVideoCodec"`
	SourceAudioCodec        string  `json:"sourceAudioCodec"`
	VideoDecision           string  `json:"videoDecision"`
	AudioDecision           string  `json:"audioDecision"`
	SubtitleDecision        string  `json:"subtitleDecision"`
	Protocol                string  `json:"protocol"`
	Container               string  `json:"container"`
	VideoCodec              string  `json:"videoCodec"`
	AudioCodec              string  `json:"audioCodec"`
	AudioChannels           int     `json:"audioChannels"`
	TranscodeHwRequested    bool    `json:"transcodeHwRequested"`
	TranscodeHwFullPipeline bool    `json:"transcodeHwFullPipeline"`
	TimeStamp               float64 `json:"timeStamp"`
}

// GetTitle returns the title of the movie, tv episode being played.  For movies, this is just the title.
// For TV Shows, it returns the show, season & episode title.
func (s Session) GetTitle() string {
	if s.Type == "episode" {
		return fmt.Sprintf("%s - S%02dE%02d - %s", s.GrandparentTitle, s.ParentIndex, s.Index, s.Title)
	}
	return s.Title

}

// GetProgress returns the progress of the session, i.e. how much of the movie / tv episode has been watched.
// Returns a percentage between 0.0 and 1.0
func (s Session) GetProgress() float64 {
	return float64(s.ViewOffset) / float64(s.Duration)
}

// GetVideoMode returns the session's video mode (transcoding, direct play, etc).
func (s Session) GetVideoMode() string {
	decisions := set.Create[string]()
	for _, media := range s.Media {
		for _, part := range media.Part {
			videoDecision := part.Decision
			if videoDecision == "transcode" {
				videoDecision = s.TranscodeSession.VideoDecision
			}
			if videoDecision == "" {
				videoDecision = "unknown"
			}
			decisions.Add(videoDecision)
		}
	}
	modes := decisions.List()
	if len(modes) == 0 {
		return "unknown"
	}
	sort.Strings(modes)
	return strings.Join(modes, ",")
}
