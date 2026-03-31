package todo

import (
	"database/sql"
	"time"
)

type Task struct {
	User_id   int       `json:"user_id"`
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}
type Mydb struct {
	db *sql.DB
}
