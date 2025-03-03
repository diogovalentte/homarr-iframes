package media

import "time"

// Calendar is a struct that represents the media releases of a specific time range
type Calendar struct {
	Releases []MediaRelease
}

// MediaRelease is a struct that represents a movie/episode on the calendar
type MediaRelease struct {
	// ReleaseDate is the date the media is released. It should have the local timezone.
	// Used for sorting and to display in the iFrame.
	ReleaseDate time.Time
	Title       string
	// Slug is usually used to generate the URL of the media in the source (*arr)
	Slug string
	// Source is a string that can be:
	// - Radarr
	// - Sonarr
	// - Lidarr
	Source         string
	PosterImageURL string
	CoverImageURL  string
	IsDownloaded   bool
	// A media should be downloaded when the its release date is after now.
	ShouldBeDownloaded bool
	// Sonnar specific
	EpisodeDetails struct {
		EpisodeName   string
		SeasonNumber  int
		EpisodeNumber int
	}
	// Lidarr specific
	ArtistDetails struct {
		ArtistName string
		Slug       string
	}
	AlbumType       string
	TrackFileCount  int
	TotalTrackCount int
}
