package plex

import "time"

// Identity contains the response of Plex's /identity API
type Identity struct {
	Size              int    `json:"size"`
	Claimed           bool   `json:"claimed"`
	MachineIdentifier string `json:"machineIdentifier"`
	Version           string `json:"version"`
}

// Sessions contains the response of Plex's /status/sessions API
type Sessions struct {
	Size     int       `json:"size"`
	Metadata []Session `json:"Metadata"`
}

// Session contains one record in a Sessions
type Session struct {
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
	User             SessionUser       `json:"User"`
	Player           SessionPlayer     `json:"Player"`
	Session          SessionStats      `json:"Session"`
	TranscodeSession SessionTranscoder `json:"TranscodeSession"`
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

type LibrariesDirectory struct {
	AllowSync        bool   `json:"allowSync"`
	Art              string `json:"art"`
	Composite        string `json:"composite"`
	Filters          bool   `json:"filters"`
	Refreshing       bool   `json:"refreshing"`
	Thumb            string `json:"thumb"`
	Key              string `json:"key"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	Agent            string `json:"agent"`
	Scanner          string `json:"scanner"`
	Language         string `json:"language"`
	Uuid             string `json:"uuid"`
	UpdatedAt        int    `json:"updatedAt"`
	CreatedAt        int    `json:"createdAt"`
	ScannedAt        int    `json:"scannedAt"`
	Content          bool   `json:"content"`
	Directory        bool   `json:"directory"`
	ContentChangedAt int    `json:"contentChangedAt"`
	Hidden           int    `json:"hidden"`
	Location         []struct {
		Id   int    `json:"id"`
		Path string `json:"path"`
	} `json:"Location"`
}

type Libraries struct {
	Size      int                  `json:"size"`
	AllowSync bool                 `json:"allowSync"`
	Title1    string               `json:"title1"`
	Directory []LibrariesDirectory `json:"Directory"`
}

type MovieLibrary struct {
	Size                int                 `json:"size"`
	AllowSync           bool                `json:"allowSync"`
	Art                 string              `json:"art"`
	Identifier          string              `json:"identifier"`
	LibrarySectionID    int                 `json:"librarySectionID"`
	LibrarySectionTitle string              `json:"librarySectionTitle"`
	LibrarySectionUUID  string              `json:"librarySectionUUID"`
	MediaTagPrefix      string              `json:"mediaTagPrefix"`
	MediaTagVersion     int                 `json:"mediaTagVersion"`
	Thumb               string              `json:"thumb"`
	Title1              string              `json:"title1"`
	Title2              string              `json:"title2"`
	ViewGroup           string              `json:"viewGroup"`
	ViewMode            int                 `json:"viewMode"`
	Metadata            []MovieLibraryEntry `json:"Metadata"`
}

type MovieLibraryEntry struct {
	RatingKey             string  `json:"ratingKey"`
	Key                   string  `json:"key"`
	Guid                  string  `json:"guid"`
	Studio                string  `json:"studio,omitempty"`
	Type                  string  `json:"type"`
	Title                 string  `json:"title"`
	ContentRating         string  `json:"contentRating,omitempty"`
	Summary               string  `json:"summary"`
	Rating                float64 `json:"rating,omitempty"`
	AudienceRating        float64 `json:"audienceRating,omitempty"`
	SkipCount             int     `json:"skipCount,omitempty"`
	LastViewedAt          int     `json:"lastViewedAt,omitempty"`
	Year                  int     `json:"year,omitempty"`
	Tagline               string  `json:"tagline,omitempty"`
	Thumb                 string  `json:"thumb"`
	Art                   string  `json:"art,omitempty"`
	Duration              int     `json:"duration"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	AddedAt               int     `json:"addedAt"`
	UpdatedAt             int     `json:"updatedAt,omitempty"`
	AudienceRatingImage   string  `json:"audienceRatingImage,omitempty"`
	PrimaryExtraKey       string  `json:"primaryExtraKey,omitempty"`
	RatingImage           string  `json:"ratingImage,omitempty"`
	Media                 []struct {
		Id              int     `json:"id"`
		Duration        int     `json:"duration"`
		Bitrate         int     `json:"bitrate"`
		Width           int     `json:"width"`
		Height          int     `json:"height"`
		AspectRatio     float64 `json:"aspectRatio"`
		AudioChannels   int     `json:"audioChannels"`
		AudioCodec      string  `json:"audioCodec"`
		VideoCodec      string  `json:"videoCodec"`
		VideoResolution string  `json:"videoResolution"`
		Container       string  `json:"container"`
		VideoFrameRate  string  `json:"videoFrameRate"`
		VideoProfile    string  `json:"videoProfile"`
		Part            []struct {
			Id                    int    `json:"id"`
			Key                   string `json:"key"`
			Duration              int    `json:"duration"`
			File                  string `json:"file"`
			Size                  int64  `json:"size"`
			Container             string `json:"container"`
			VideoProfile          string `json:"videoProfile"`
			AudioProfile          string `json:"audioProfile,omitempty"`
			Has64BitOffsets       bool   `json:"has64bitOffsets,omitempty"`
			OptimizedForStreaming bool   `json:"optimizedForStreaming,omitempty"`
			HasThumbnail          string `json:"hasThumbnail,omitempty"`
		} `json:"Part"`
		OptimizedForStreaming int    `json:"optimizedForStreaming,omitempty"`
		AudioProfile          string `json:"audioProfile,omitempty"`
		Has64BitOffsets       bool   `json:"has64bitOffsets,omitempty"`
	} `json:"Media"`
	Genre []struct {
		Tag string `json:"tag"`
	} `json:"Genre,omitempty"`
	Director []struct {
		Tag string `json:"tag"`
	} `json:"Director,omitempty"`
	Writer []struct {
		Tag string `json:"tag"`
	} `json:"Writer,omitempty"`
	Country []struct {
		Tag string `json:"tag"`
	} `json:"Country,omitempty"`
	Role []struct {
		Tag string `json:"tag"`
	} `json:"Role,omitempty"`
	ViewCount  int `json:"viewCount,omitempty"`
	Collection []struct {
		Tag string `json:"tag"`
	} `json:"Collection,omitempty"`
	ChapterSource string  `json:"chapterSource,omitempty"`
	TitleSort     string  `json:"titleSort,omitempty"`
	OriginalTitle string  `json:"originalTitle,omitempty"`
	UserRating    float64 `json:"userRating,omitempty"`
	LastRatedAt   int     `json:"lastRatedAt,omitempty"`
}

type ShowLibrary struct {
	Size                int                `json:"size"`
	AllowSync           bool               `json:"allowSync"`
	Art                 string             `json:"art"`
	Identifier          string             `json:"identifier"`
	LibrarySectionID    int                `json:"librarySectionID"`
	LibrarySectionTitle string             `json:"librarySectionTitle"`
	LibrarySectionUUID  string             `json:"librarySectionUUID"`
	MediaTagPrefix      string             `json:"mediaTagPrefix"`
	MediaTagVersion     int                `json:"mediaTagVersion"`
	Nocache             bool               `json:"nocache"`
	Thumb               string             `json:"thumb"`
	Title1              string             `json:"title1"`
	Title2              string             `json:"title2"`
	ViewGroup           string             `json:"viewGroup"`
	ViewMode            int                `json:"viewMode"`
	Metadata            []ShowLibraryEntry `json:"Metadata"`
}

type ShowLibraryEntry struct {
	RatingKey             string  `json:"ratingKey"`
	Key                   string  `json:"key"`
	Guid                  string  `json:"guid"`
	Studio                string  `json:"studio"`
	Type                  string  `json:"type"`
	Title                 string  `json:"title"`
	ContentRating         string  `json:"contentRating,omitempty"`
	Summary               string  `json:"summary"`
	Index                 int     `json:"index"`
	AudienceRating        float64 `json:"audienceRating,omitempty"`
	ViewCount             int     `json:"viewCount,omitempty"`
	SkipCount             int     `json:"skipCount,omitempty"`
	LastViewedAt          int     `json:"lastViewedAt,omitempty"`
	Year                  int     `json:"year"`
	Thumb                 string  `json:"thumb"`
	Art                   string  `json:"art"`
	Theme                 string  `json:"theme,omitempty"`
	Duration              int     `json:"duration"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	LeafCount             int     `json:"leafCount"`
	ViewedLeafCount       int     `json:"viewedLeafCount"`
	ChildCount            int     `json:"childCount"`
	AddedAt               int     `json:"addedAt"`
	UpdatedAt             int     `json:"updatedAt"`
	AudienceRatingImage   string  `json:"audienceRatingImage,omitempty"`
	PrimaryExtraKey       string  `json:"primaryExtraKey,omitempty"`
	Genre                 []struct {
		Tag string `json:"tag"`
	} `json:"Genre"`
	Country []struct {
		Tag string `json:"tag"`
	} `json:"Country,omitempty"`
	Role []struct {
		Tag string `json:"tag"`
	} `json:"Role,omitempty"`
	Tagline       string  `json:"tagline,omitempty"`
	TitleSort     string  `json:"titleSort,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	Banner        string  `json:"banner,omitempty"`
	OriginalTitle string  `json:"originalTitle,omitempty"`
}

type AccessToken struct {
	Type      string    `json:"type"`
	Device    string    `json:"device,omitempty"`
	Token     string    `json:"token"`
	Owned     bool      `json:"owned"`
	CreatedAt time.Time `json:"createdAt"`
	Invited   struct {
		Id       int         `json:"id"`
		Uuid     string      `json:"uuid"`
		Title    string      `json:"title"`
		Username interface{} `json:"username"`
		Thumb    string      `json:"thumb"`
		Profile  struct {
			AutoSelectAudio              bool        `json:"autoSelectAudio"`
			DefaultAudioLanguage         interface{} `json:"defaultAudioLanguage"`
			DefaultSubtitleLanguage      interface{} `json:"defaultSubtitleLanguage"`
			AutoSelectSubtitle           int         `json:"autoSelectSubtitle"`
			DefaultSubtitleAccessibility int         `json:"defaultSubtitleAccessibility"`
			DefaultSubtitleForced        int         `json:"defaultSubtitleForced"`
		} `json:"profile"`
		Scrobbling    []interface{} `json:"scrobbling"`
		ScrobbleTypes string        `json:"scrobbleTypes"`
	} `json:"invited,omitempty"`
	Settings struct {
		AllowChannels      bool        `json:"allowChannels"`
		FilterMovies       *string     `json:"filterMovies"`
		FilterMusic        *string     `json:"filterMusic"`
		FilterPhotos       interface{} `json:"filterPhotos"`
		FilterTelevision   *string     `json:"filterTelevision"`
		FilterAll          interface{} `json:"filterAll"`
		AllowSync          bool        `json:"allowSync"`
		AllowCameraUpload  bool        `json:"allowCameraUpload"`
		AllowSubtitleAdmin bool        `json:"allowSubtitleAdmin"`
		AllowTuners        int         `json:"allowTuners"`
	} `json:"settings,omitempty"`
	Sections []struct {
		Key       int       `json:"key"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"sections,omitempty"`
}
