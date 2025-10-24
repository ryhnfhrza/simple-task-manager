package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
	"github.com/ryhnfhrza/simple-task-manager/service"
)

type userControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userControllerImpl{
		UserService: userService,
	}
}

func (controller *userControllerImpl) Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userRequest := web.UserRegisterRequest{}
	helper.ReadFromRequestBody(request, &userRequest)

	userResponse, err := controller.UserService.Register(request.Context(), userRequest)
	helper.PanicIfError(err)

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *userControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userRequest := web.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &userRequest)

	userResponse, err := controller.UserService.Login(request.Context(), userRequest)
	helper.PanicIfError(err)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
