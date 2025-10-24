package service

import (
	"context"

	"github.com/ryhnfhrza/simple-task-manager/model/web"
)

type UserService interface {
	Register(ctx context.Context, req web.UserRegisterRequest) (web.UserResponse, error)
	Login(ctx context.Context, req web.UserLoginRequest) (web.UserLoginResponse, error)
}
