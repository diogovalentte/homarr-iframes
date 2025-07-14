package jellyfin

// Item represents a single Jellyfin media item
type Item struct {
	Name              string            `json:"Name"`
	ServerId          string            `json:"ServerId"`
	ID                string            `json:"Id"`
	Year              int               `json:"ProductionYear"`
	Type              string            `json:"Type"`
	ImageTags         map[string]string `json:"ImageTags"`
	BackdropImageTags []string          `json:"BackdropImageTags"`
	PrimaryImageURL   string
	BackdropImageURL  string
	ItemURL           string
}
