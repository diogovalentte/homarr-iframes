package media

// Calendar is a struct that represents the media releases of a specific time range
type Calendar struct {
	Releases []MediaRelease
}

// MediaRelease is a struct that represents a movie/episode on the calendar
type MediaRelease struct {
	Title string
	// Slug is usually used to generate the URL of the media in the source (*arr)
	Slug string
	// Source is a string that can be:
	// - Radarr
	// - Sonarr
	Source         string
	CoverImageURL  string
	IsDownloaded   bool
	EpisodeDetails struct {
		SeasonNumber  int
		EpisodeNumber int
		EpisodeName   string
	}
}

type defaultReleaseImagesResponse struct {
	CoverType string `json:"coverType"`
	RemoteURL string `json:"remoteUrl"`
}
