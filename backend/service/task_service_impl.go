package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ryhnfhrza/simple-task-manager/exception"
	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/model/domain"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
	"github.com/ryhnfhrza/simple-task-manager/repository"
)

type taskServiceImpl struct {
	TaskRepository repository.TaskRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewTaskService(taskRepository repository.TaskRepository, db *sql.DB, validate *validator.Validate) TaskService {
	return &taskServiceImpl{
		TaskRepository: taskRepository,
		DB:             db,
		Validate:       validate,
	}
}

func (service *taskServiceImpl) Create(ctx context.Context, req web.TaskCreateRequest) web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.Validate.Struct(req)
	helper.PanicIfError(err)

	if req.Title == "" {
		req.Title = "No title"
	}
	if req.Description == "" {
		req.Description = "No Desc"
	}

	task := domain.Task{
		UserId:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		Completed:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.DueDate != nil {
		task.DueDate = sql.NullTime{Time: *req.DueDate, Valid: true}
	} else {
		task.DueDate = sql.NullTime{Valid: false}
	}

	task = service.TaskRepository.SaveTask(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

func (service *taskServiceImpl) Update(ctx context.Context, req web.TaskUpdateRequest) web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.Validate.Struct(req)
	helper.PanicIfError(err)

	task, err := service.TaskRepository.FindTaskById(ctx, tx, req.Id, req.UserId)
	if err != nil {
		panic(exception.NewNotFoundError("task not found"))
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.DueDate != nil && !req.DueDate.IsZero() {
		task.DueDate = sql.NullTime{Time: *req.DueDate, Valid: true}
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}
	task.UserId = int64(req.UserId)

	task = service.TaskRepository.UpdateTask(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

func (service *taskServiceImpl) Delete(ctx context.Context, taskId, userId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task, err := service.TaskRepository.FindTaskById(ctx, tx, taskId, userId)
	if err != nil {
		panic(exception.NewNotFoundError("task not found"))
	}
	task.UserId = int64(userId)
	service.TaskRepository.DeleteTask(ctx, tx, task)
}

func (service *taskServiceImpl) FindById(ctx context.Context, taskId, userId int) web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task, err := service.TaskRepository.FindTaskById(ctx, tx, taskId, userId)
	if err != nil {
		panic(exception.NewNotFoundError("task not found"))
	}

	return helper.ToTaskResponse(task)
}

func (service *taskServiceImpl) FindAll(ctx context.Context, userId int64, filter domain.TaskFilter) []web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.Order == "" {
		filter.Order = "DESC"
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	tasks := service.TaskRepository.FindAllTask(ctx, tx, userId, filter)

	return helper.ToTaskResponses(tasks)
}
