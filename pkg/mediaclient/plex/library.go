package plex

import (
	"context"
	"fmt"
)

type Libraries struct {
	Size      int                  `json:"size"`
	AllowSync bool                 `json:"allowSync"`
	Title1    string               `json:"title1"`
	Directory []LibrariesDirectory `json:"Directory"`
}

type LibrariesDirectory struct {
	AllowSync        bool      `json:"allowSync"`
	Art              string    `json:"art"`
	Composite        string    `json:"composite"`
	Filters          bool      `json:"filters"`
	Refreshing       bool      `json:"refreshing"`
	Thumb            string    `json:"thumb"`
	Key              string    `json:"key"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Agent            string    `json:"agent"`
	Scanner          string    `json:"scanner"`
	Language         string    `json:"language"`
	Uuid             string    `json:"uuid"`
	UpdatedAt        Timestamp `json:"updatedAt"`
	CreatedAt        Timestamp `json:"createdAt"`
	ScannedAt        Timestamp `json:"scannedAt"`
	Content          bool      `json:"content"`
	Directory        bool      `json:"directory"`
	ContentChangedAt int       `json:"contentChangedAt"`
	Hidden           int       `json:"hidden"`
	Location         []struct {
		Id   int    `json:"id"`
		Path string `json:"path"`
	} `json:"Location"`
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
	RatingKey             string    `json:"ratingKey"`
	Key                   string    `json:"key"`
	Guid                  string    `json:"guid"`
	Studio                string    `json:"studio,omitempty"`
	Type                  string    `json:"type"`
	Title                 string    `json:"title"`
	ContentRating         string    `json:"contentRating,omitempty"`
	Summary               string    `json:"summary"`
	Rating                float64   `json:"rating,omitempty"`
	AudienceRating        float64   `json:"audienceRating,omitempty"`
	SkipCount             int       `json:"skipCount,omitempty"`
	LastViewedAt          Timestamp `json:"lastViewedAt,omitempty"`
	Year                  int       `json:"year,omitempty"`
	Tagline               string    `json:"tagline,omitempty"`
	Thumb                 string    `json:"thumb"`
	Art                   string    `json:"art,omitempty"`
	Duration              int       `json:"duration"`
	OriginallyAvailableAt string    `json:"originallyAvailableAt"`
	AddedAt               Timestamp `json:"addedAt"`
	UpdatedAt             Timestamp `json:"updatedAt,omitempty"`
	AudienceRatingImage   string    `json:"audienceRatingImage,omitempty"`
	PrimaryExtraKey       string    `json:"primaryExtraKey,omitempty"`
	RatingImage           string    `json:"ratingImage,omitempty"`
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
	RatingKey             string    `json:"ratingKey"`
	Key                   string    `json:"key"`
	Guid                  string    `json:"guid"`
	Studio                string    `json:"studio"`
	Type                  string    `json:"type"`
	Title                 string    `json:"title"`
	ContentRating         string    `json:"contentRating,omitempty"`
	Summary               string    `json:"summary"`
	Index                 int       `json:"index"`
	AudienceRating        float64   `json:"audienceRating,omitempty"`
	ViewCount             int       `json:"viewCount,omitempty"`
	SkipCount             int       `json:"skipCount,omitempty"`
	LastViewedAt          Timestamp `json:"lastViewedAt,omitempty"`
	Year                  int       `json:"year"`
	Thumb                 string    `json:"thumb"`
	Art                   string    `json:"art"`
	Theme                 string    `json:"theme,omitempty"`
	Duration              int       `json:"duration"`
	OriginallyAvailableAt string    `json:"originallyAvailableAt"`
	LeafCount             int       `json:"leafCount"`
	ViewedLeafCount       int       `json:"viewedLeafCount"`
	ChildCount            int       `json:"childCount"`
	AddedAt               Timestamp `json:"addedAt"`
	UpdatedAt             Timestamp `json:"updatedAt"`
	AudienceRatingImage   string    `json:"audienceRatingImage,omitempty"`
	PrimaryExtraKey       string    `json:"primaryExtraKey,omitempty"`
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

func (c *Client) GetLibraries(ctx context.Context) (libraries Libraries, err error) {
	err = c.call(ctx, "/library/sections", &libraries)
	return
}

func (c *Client) GetMovieLibrary(ctx context.Context, key string) (library MovieLibrary, err error) {
	err = c.call(ctx, fmt.Sprintf("/library/sections/%s/all", key), &library)
	return
}

func (c *Client) GetShowLibrary(ctx context.Context, key string) (library ShowLibrary, err error) {
	err = c.call(ctx, fmt.Sprintf("/library/sections/%s/all", key), &library)
	return
}
