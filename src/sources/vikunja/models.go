package vikunja

import "time"

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Token string `json:"token"`
}

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created"`
	DueDate   time.Time `json:"due_date"`
	EndDate   time.Time `json:"end_date"`
	Priority  int       `json:"priority"`
}
