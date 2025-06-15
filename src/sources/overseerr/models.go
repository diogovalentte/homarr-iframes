package overseerr

type Request struct {
	RequestedBy RequestedBy `json:"requestedBy"`
	Media       Media       `json:"media"`
	ID          int         `json:"id"`
	Status      int         `json:"status"`
}

type Media struct {
	Type   string `json:"mediaType"`
	IMDBID string `json:"imdbId"`
	ID     int    `json:"id"`
	Status int    `json:"status"`
	TMDBID int    `json:"tmdbId"`
	TVDBID int    `json:"tvdbId"`
}

type RequestedBy struct {
	Username string `json:"displayName"`
	Avatar   string `json:"avatar"`
	ID       int    `json:"id"`
}

// GenericMedia is a generic media struct used for both movies and tv shows when requesting a media
type GenericMedia struct {
	Name         string
	ReleaseDate  string
	PosterPath   string
	BackdropPath string
	ID           int
}

type IframeRequestData struct {
	Status IframeStatus
	Media  struct {
		Name        string
		Type        string
		Year        string
		BackdropURL string
		PosterURL   string
		URL         string
		TMDBID      int
	}
	Request struct {
		Username       string
		AvatarURL      string
		UserProfileURL string
		UserID         int
	}
}

type IframeStatus struct {
	Status          string
	Color           string
	BackgroundColor string
}
