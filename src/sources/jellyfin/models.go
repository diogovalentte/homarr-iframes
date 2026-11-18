package jellyfin

// Item represents a single Jellyfin media item
type Item struct {
	BackdropImageTags []string          `json:"BackdropImageTags"`
	ImageTags         map[string]string `json:"ImageTags"`
	Name              string            `json:"Name"`
	ServerId          string            `json:"ServerId"`
	ID                string            `json:"Id"`
	Type              string            `json:"Type"`
	SeriesID          string            `json:"SeriesId"`
	SeriesName        string            `json:"SeriesName"`
	EpisodeURL        string
	PrimaryImageURL   string
	BackdropImageURL  string
	ItemURL           string
	Year              int   `json:"ProductionYear"`
	SeasonNumber      int   `json:"ParentIndexNumber"`
	EpisodeNumber     int   `json:"IndexNumber"`
	RunTimeTicks      int64 `json:"RunTimeTicks"`
}

// PlayState represents the current playback state of a session
type PlayState struct {
	PositionTicks int64 `json:"PositionTicks"`
	IsPaused      bool  `json:"IsPaused"`
}

// Session represents a user session in Jellyfin
type Session struct {
	PlayState      PlayState `json:"PlayState"`
	NowPlayingItem *Item     `json:"NowPlayingItem"`
	ID             string    `json:"Id"`
	UserName       string    `json:"UserName"`
	UserID         string    `json:"UserId"`
	UserAvatarURL  string
}
