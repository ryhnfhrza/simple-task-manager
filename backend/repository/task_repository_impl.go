package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/model/domain"
)

type taskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &taskRepositoryImpl{}
}

func (repo *taskRepositoryImpl) SaveTask(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	query := "insert into tasks(title,description,due_date,completed,user_id) values (?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, query, task.Title, task.Description, task.DueDate, task.Completed, task.UserId)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	task.Id = id

	return task
}

func (repo *taskRepositoryImpl) UpdateTask(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	query := "update tasks set title = ?,description = ?,due_date = ?, completed = ? where id = ? and user_id = ?"

	_, err := tx.ExecContext(ctx, query, task.Title, task.Description, task.DueDate, task.Completed, task.Id, task.UserId)
	helper.PanicIfError(err)

	return task
}

func (repo *taskRepositoryImpl) DeleteTask(ctx context.Context, tx *sql.Tx, task domain.Task) {
	query := "delete from tasks where id = ? and user_id = ?"
	_, err := tx.ExecContext(ctx, query, task.Id, task.UserId)
	helper.PanicIfError(err)

}

func (repo *taskRepositoryImpl) FindTaskById(ctx context.Context, tx *sql.Tx, taskId, userId int) (domain.Task, error) {
	query := "select id,title,description,due_date,completed,created_at,updated_at from tasks where id = ? and user_id = ?"
	row := tx.QueryRowContext(ctx, query, taskId, userId)

	task := domain.Task{}

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (repo *taskRepositoryImpl) FindAllTask(ctx context.Context, tx *sql.Tx, userId int64, filter domain.TaskFilter) []domain.Task {
	var (
		sb   strings.Builder
		args []interface{}
	)

	sb.WriteString(`
        SELECT
            id,
            title,
            description,
            due_date,
            completed,
            created_at,
            updated_at
        FROM tasks
        WHERE user_id = ?
    `)
	args = append(args, userId)

	if filter.Completed != nil {
		sb.WriteString(" AND completed = ?")
		args = append(args, *filter.Completed)
	}

	if filter.DueBefore != nil {
		sb.WriteString(" AND due_date < ?")
		args = append(args, filter.DueBefore.Format(time.RFC3339))
	}
	if filter.DueAfter != nil {
		sb.WriteString(" AND due_date > ?")
		args = append(args, filter.DueAfter.Format(time.RFC3339))
	}

	if filter.SortBy != "" {
		dir := strings.ToUpper(filter.Order)
		if dir != "ASC" && dir != "DESC" {
			dir = "ASC"
		}
		sb.WriteString(fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, dir))
	}

	if filter.Limit > 0 {
		sb.WriteString(" LIMIT ?")
		args = append(args, filter.Limit)
		if filter.Offset > 0 {
			sb.WriteString(" OFFSET ?")
			args = append(args, filter.Offset)
		}
	}

	rows, err := tx.QueryContext(ctx, sb.String(), args...)
	helper.PanicIfError(err)
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var rawComp int16

		helper.PanicIfError(rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&rawComp,
			&task.CreatedAt,
			&task.UpdatedAt,
		))

		task.Completed = rawComp
		task.UserId = userId
		tasks = append(tasks, task)
	}

	return tasks
}
