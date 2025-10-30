package web

import (
	"time"

	"github.com/ryhnfhrza/simple-task-manager/internal/types"
)

type TaskResponse struct {
	Id          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	DueDate     *types.CustomTime `json:"due_date,omitempty"`
	Completed   int16             `json:"completed"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
