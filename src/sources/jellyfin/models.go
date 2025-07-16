package jellyfin

// Item represents a single Jellyfin media item
type Item struct {
	BackdropImageTags []string          `json:"BackdropImageTags"`
	ImageTags         map[string]string `json:"ImageTags"`
	Name              string            `json:"Name"`
	ServerId          string            `json:"ServerId"`
	ID                string            `json:"Id"`
	Type              string            `json:"Type"`
	PrimaryImageURL   string
	BackdropImageURL  string
	ItemURL           string
	Year              int `json:"ProductionYear"`
}
