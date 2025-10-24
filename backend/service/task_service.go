package service

import (
	"context"

	"github.com/ryhnfhrza/simple-task-manager/model/domain"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
)

type TaskService interface {
	Create(ctx context.Context, req web.TaskCreateRequest) web.TaskResponse
	Update(ctx context.Context, req web.TaskUpdateRequest) web.TaskResponse
	Delete(ctx context.Context, taskId, userId int)
	FindById(ctx context.Context, taskId, userId int) web.TaskResponse
	FindAll(ctx context.Context, userId int64, filter domain.TaskFilter) []web.TaskResponse
}
