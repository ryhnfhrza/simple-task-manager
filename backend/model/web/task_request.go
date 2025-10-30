package web

import (
	"github.com/ryhnfhrza/simple-task-manager/internal/types"
)

type TaskCreateRequest struct {
	UserId      int64             `json:"user_id" validate:"required"`
	Title       string            `json:"title" validate:"required_without=Description"`
	Description string            `json:"description" validate:"required_without=Title"`
	DueDate     *types.CustomTime `json:"due_date"`
}

type TaskUpdateRequest struct {
	Id          int               `json:"id" validate:"required"`
	UserId      int               `json:"user_id" validate:"required"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Completed   *int16            `json:"completed,omitempty" validate:"omitempty,oneof=0 1"`
	DueDate     *types.CustomTime `json:"due_date,omitempty"`
}
