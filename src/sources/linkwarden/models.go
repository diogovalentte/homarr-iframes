package linkwarden

import "time"

// Link represents a Linkwarden link
type Link struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Description   *string     `json:"description"`
	CollectionID  *int        `json:"collectionId"`
	URL           string      `json:"url"`
	TextContent   *string     `json:"textContent"`
	Preview       string      `json:"preview"`
	ImagePath     *string     `json:"image"`
	PDFPath       *string     `json:"pdf"`
	ReadablePath  *string     `json:"readable"`
	LastPreserved time.Time   `json:"lastPreserved"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	Tags          []*Tag      `json:"tags"`
	Collection    *Collection `json:"collection"`
	PinnedBy      []*string   `json:"pinnedBy"`
}

// Collection represents a Linkwarden collection
type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Color       string    `json:"color"`
	ParentID    *int      `json:"parentId,omitempty`
	IsPublic    bool      `json:"isPublic"`
	OwnerID     *int      `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Tag represents a Linkwarden tag
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   *int      `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
