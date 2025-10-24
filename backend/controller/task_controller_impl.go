package controller

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/ryhnfhrza/simple-task-manager/exception"
	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/model/domain"
	"github.com/ryhnfhrza/simple-task-manager/model/web"
	"github.com/ryhnfhrza/simple-task-manager/service"
)

type taskControllerImpl struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return &taskControllerImpl{
		TaskService: taskService,
	}
}

func (controller *taskControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskCreateRequest := web.TaskCreateRequest{}
	helper.ReadFromRequestBody(request, &taskCreateRequest)

	userId, ok := helper.GetUserIDFromContext(request.Context())
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskCreateRequest.UserId = userId

	taskResponse := controller.TaskService.Create(request.Context(), taskCreateRequest)

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *taskControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskUpdateRequest := web.TaskUpdateRequest{}
	helper.ReadFromRequestBody(request, &taskUpdateRequest)

	userId, ok := helper.GetUserIDFromContext(request.Context())
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskIdString := params.ByName("taskId")
	taskId, err := strconv.Atoi(taskIdString)
	helper.PanicIfError(err)

	taskUpdateRequest.Id = taskId
	taskUpdateRequest.UserId = int(userId)

	taskResponse := controller.TaskService.Update(request.Context(), taskUpdateRequest)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *taskControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId, ok := helper.GetUserIDFromContext(request.Context())
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskIdString := params.ByName("taskId")
	taskId, err := strconv.Atoi(taskIdString)
	helper.PanicIfError(err)

	controller.TaskService.Delete(request.Context(), taskId, int(userId))

	webResponse := web.WebResponse{
		Code:   http.StatusNoContent,
		Status: "NO CONTENT",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *taskControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId, ok := helper.GetUserIDFromContext(request.Context())
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	taskIdString := params.ByName("taskId")
	taskId, err := strconv.Atoi(taskIdString)
	helper.PanicIfError(err)

	taskResponse := controller.TaskService.FindById(request.Context(), taskId, int(userId))

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *taskControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId, ok := helper.GetUserIDFromContext(request.Context())
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}
	query := request.URL.Query()

	var filter domain.TaskFilter

	if completedStr := query.Get("completed"); completedStr != "" {
		c, err := strconv.Atoi(completedStr)
		if err == nil {
			filter.Completed = &c
		}
	}

	if dueBeforeStr := query.Get("due_before"); dueBeforeStr != "" {
		parsed, err := helper.ParseFlexibleDate(dueBeforeStr)
		if err != nil {
			panic(exception.NewBadRequest("Invalid due_before format. Use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS"))
		}
		filter.DueBefore = &parsed
	}
	if dueAfterStr := query.Get("due_after"); dueAfterStr != "" {
		parsed, err := helper.ParseFlexibleDate(dueAfterStr)
		if err != nil {
			panic(exception.NewBadRequest("Invalid due_after format. Use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS"))
		}
		filter.DueAfter = &parsed
	}

	filter.SortBy = query.Get("sort_by")
	filter.Order = query.Get("order")

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = l
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = o
		}
	}

	taskResponse := controller.TaskService.FindAll(request.Context(), int64(userId), filter)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
