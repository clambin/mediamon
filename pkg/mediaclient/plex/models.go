package plex

// IdentityResponse contains the response of Plex' /identity endpoint
type IdentityResponse struct {
	MediaContainer struct {
		Size              int    `json:"size"`
		Claimed           bool   `json:"claimed"`
		MachineIdentifier string `json:"machineIdentifier"`
		Version           string `json:"version"`
	} `json:"MediaContainer"`
}

// SessionsResponse contains the response of Plex' /status/sessions endpoint
type SessionsResponse struct {
	MediaContainer struct {
		Size     int                      `json:"size"`
		Metadata []SessionsResponseRecord `json:"Metadata"`
	} `json:"MediaContainer"`
}

// SessionsResponseRecord contains one record in a SessionsResponse
type SessionsResponseRecord struct {
	AddedAt               int     `json:"addedAt"`
	Art                   string  `json:"art"`
	AudienceRating        float64 `json:"audienceRating"`
	AudienceRatingImage   string  `json:"audienceRatingImage"`
	ContentRating         string  `json:"contentRating"`
	Duration              int     `json:"duration"`
	GrandparentArt        string  `json:"grandparentArt"`
	GrandparentGUID       string  `json:"grandparentGuid"`
	GrandparentKey        string  `json:"grandparentKey"`
	GrandparentRatingKey  string  `json:"grandparentRatingKey"`
	GrandparentTheme      string  `json:"grandparentTheme"`
	GrandparentThumb      string  `json:"grandparentThumb"`
	GrandparentTitle      string  `json:"grandparentTitle"`
	GUID                  string  `json:"guid"`
	Index                 int     `json:"index"`
	Key                   string  `json:"key"`
	LastViewedAt          int     `json:"lastViewedAt"`
	LibrarySectionID      string  `json:"librarySectionID"`
	LibrarySectionKey     string  `json:"librarySectionKey"`
	LibrarySectionTitle   string  `json:"librarySectionTitle"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	ParentGUID            string  `json:"parentGuid"`
	ParentIndex           int     `json:"parentIndex"`
	ParentKey             string  `json:"parentKey"`
	ParentRatingKey       string  `json:"parentRatingKey"`
	ParentThumb           string  `json:"parentThumb"`
	ParentTitle           string  `json:"parentTitle"`
	Rating                float64 `json:"rating"`
	RatingKey             string  `json:"ratingKey"`
	SessionKey            string  `json:"sessionKey"`
	Summary               string  `json:"summary"`
	Thumb                 string  `json:"thumb"`
	Title                 string  `json:"title"`
	Type                  string  `json:"type"`
	UpdatedAt             int     `json:"updatedAt"`
	ViewOffset            int     `json:"viewOffset"`
	Media                 []struct {
		AudioProfile          string `json:"audioProfile"`
		ID                    string `json:"id"`
		VideoProfile          string `json:"videoProfile"`
		AudioChannels         int    `json:"audioChannels"`
		AudioCodec            string `json:"audioCodec"`
		Bitrate               int    `json:"bitrate"`
		Container             string `json:"container"`
		Duration              int    `json:"duration"`
		Height                int    `json:"height"`
		OptimizedForStreaming bool   `json:"optimizedForStreaming"`
		Protocol              string `json:"protocol"`
		VideoCodec            string `json:"videoCodec"`
		VideoFrameRate        string `json:"videoFrameRate"`
		VideoResolution       string `json:"videoResolution"`
		Width                 int    `json:"width"`
		Selected              bool   `json:"selected"`
		Part                  []struct {
			AudioProfile          string `json:"audioProfile"`
			ID                    string `json:"id"`
			VideoProfile          string `json:"videoProfile"`
			Bitrate               int    `json:"bitrate"`
			Container             string `json:"container"`
			Duration              int    `json:"duration"`
			Height                int    `json:"height"`
			OptimizedForStreaming bool   `json:"optimizedForStreaming"`
			Protocol              string `json:"protocol"`
			Width                 int    `json:"width"`
			Decision              string `json:"decision"`
			Selected              bool   `json:"selected"`
			Stream                []struct {
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
			} `json:"Stream"`
		} `json:"Part"`
	} `json:"Media"`
	Director []struct {
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
	User             SessionsResponseRecordUser             `json:"User"`
	Player           SessionsResponseRecordPlayer           `json:"Player"`
	Session          SessionsResponseRecordSession          `json:"Session"`
	TranscodeSession SessionsResponseRecordTranscodeSession `json:"TranscodeSession"`
}

// SessionsResponseRecordUser contains the user details inside a SessionsResponseRecord
type SessionsResponseRecordUser struct {
	ID    string `json:"id"`
	Thumb string `json:"thumb"`
	Title string `json:"title"`
}

// SessionsResponseRecordPlayer contains the player details inside a SessionsResponseRecord
type SessionsResponseRecordPlayer struct {
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

// SessionsResponseRecordSession contains the session details inside a SessionsResponseRecord
type SessionsResponseRecordSession struct {
	ID        string `json:"id"`
	Bandwidth int    `json:"bandwidth"`
	Location  string `json:"location"`
}

// SessionsResponseRecordTranscodeSession contains the transcoder details inside a SessionsResponseRecord.
// If the session doesn't transcode any media streams, all fields will be blank.
type SessionsResponseRecordTranscodeSession struct {
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
