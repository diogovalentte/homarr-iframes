package linkwarden

import "time"

// Link represents a Linkwarden link
// ! IMPORTANT !
// If you add a filed where the value is a pointer,
// you have to update the Linkwarden.GetHash method
// to set it to nil.
type Link struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	CollectionID *int        `json:"collectionId"`
	URL          string      `json:"url"`
	CreatedAt    time.Time   `json:"createdAt"`
	Collection   *Collection `json:"collection"`
}

// Collection represents a Linkwarden collection
type Collection struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
