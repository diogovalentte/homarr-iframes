package vikunja

import "time"

type Info struct {
	Version string `json:"version"`
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Token string `json:"token"`
}

// Task represents a Vikunja task
// ! IMPORTANT !
// If you add a filed where the value is a pointer,
// you have to update the Vikunja.GetHash method
// to set it to nil.
type Task struct {
	ID          int       `json:"id"`
	Done        bool      `json:"done"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created"`
	DueDate     time.Time `json:"due_date"`
	EndDate     time.Time `json:"end_date"`
	Priority    int       `json:"priority"`
	RepeatAfter int       `json:"repeat_after"`
	RepeatMode  int       `json:"repeat_mode"`
	ProjectID   int       `json:"project_id"`
	IsFavorite  bool      `json:"is_favorite"`
}

// Project represents a Vikunja project
// ! IMPORTANT !
// If you add a filed where the value is a pointer,
// you have to update the Vikunja.GetHash method
// to set it to nil.
type Project struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	HexColor string `json:"hex_color"`
}
