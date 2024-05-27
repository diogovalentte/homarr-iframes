package overseerr

type Request struct {
	ID          int         `json:"id"`
	Status      int         `json:"status"`
	Media       Media       `json:"media"`
	RequestedBy RequestedBy `json:"requestedBy"`
}

type Media struct {
	ID     int    `json:"id"`
	Type   string `json:"mediaType"`
	Status int    `json:"status"`
	TMDBID int    `json:"tmdbId"`
	TVDBID int    `json:"tvdbId"`
	IMDBID int    `json:"imdbId"`
}

type RequestedBy struct {
	ID       int    `json:"id"`
	Username string `json:"displayName"`
	Avatar   string `json:"avatar"`
}

// GenericMedia is a generic media struct used for both movies and tv shows when requesting a media
type GenericMedia struct {
	Name        string
	ID          int
	ReleaseDate string
	PosterPath  string
}
