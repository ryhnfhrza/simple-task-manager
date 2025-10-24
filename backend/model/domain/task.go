package domain

import (
	"database/sql"
	"time"
)

type Task struct {
	Id          int64
	Title       string
	Description string
	DueDate     sql.NullTime
	Completed   int16
	UserId      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
