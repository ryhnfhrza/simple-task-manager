package repository

import (
	"context"
	"database/sql"

	"github.com/ryhnfhrza/simple-task-manager/model/domain"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task
	UpdateTask(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task
	DeleteTask(ctx context.Context, tx *sql.Tx, task domain.Task)
	FindTaskById(ctx context.Context, tx *sql.Tx, taskId, userId int) (domain.Task, error)
	FindAllTask(ctx context.Context, tx *sql.Tx, userId int64, filter domain.TaskFilter) []domain.Task
}
