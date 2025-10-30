package helper

import (
	"github.com/ryhnfhrza/simple-task-manager/internal/types"
	"github.com/ryhnfhrza/simple-task-manager/model/domain"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{

		Username: user.Username,
	}
}

func ToUserLoginResponse(user domain.User, token string) web.UserLoginResponse {
	return web.UserLoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    token,
	}
}

func ToTaskResponse(task domain.Task) web.TaskResponse {
	var dueDate *types.CustomTime
	if task.DueDate.Valid {
		dueDate = &types.CustomTime{Time: task.DueDate.Time}
	}

	return web.TaskResponse{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     dueDate,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ToTaskResponses(tasks []domain.Task) []web.TaskResponse {
	var taskResponses []web.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, ToTaskResponse(task))
	}

	return taskResponses
}
