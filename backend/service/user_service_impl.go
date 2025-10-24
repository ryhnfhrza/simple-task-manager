package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/ryhnfhrza/simple-task-manager/exception"
	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/util"

	"github.com/ryhnfhrza/simple-task-manager/model/domain"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
	"github.com/ryhnfhrza/simple-task-manager/repository"
)

type userServiceImpl struct {
	UserRepository repository.UsersRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UsersRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &userServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *userServiceImpl) Register(ctx context.Context, req web.UserRegisterRequest) (web.UserResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.Validate.Struct(req)
	helper.PanicIfError(err)

	_, err = service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err == nil {
		panic(exception.NewConflictError("username already exists"))
	}

	hashPassword, err := util.HashPassword(req.Password)
	helper.PanicIfError(err)

	user := domain.User{
		Username:     req.Username,
		PasswordHash: hashPassword,
	}

	err = service.UserRepository.Save(ctx, tx, &user)
	helper.PanicIfError(err)

	return helper.ToUserResponse(user), nil
}

func (service *userServiceImpl) Login(ctx context.Context, req web.UserLoginRequest) (web.UserLoginResponse, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.Validate.Struct(req)
	helper.PanicIfError(err)

	user, err := service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		panic(exception.NewUnauthorizedError("invalid username or password"))
	}

	if !util.CheckPasswordHash(req.Password, user.PasswordHash) {
		panic(exception.NewUnauthorizedError("invalid username or password"))
	}

	token, err := util.CreateToken(user)
	helper.PanicIfError(err)

	return helper.ToUserLoginResponse(*user, token), nil
}
