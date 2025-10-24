package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ryhnfhrza/simple-task-manager/controller"
	"github.com/ryhnfhrza/simple-task-manager/middleware"
)

func NewRouter(userController controller.UserController, taskController controller.TaskController) *httprouter.Router {
	router := httprouter.New()

	//user router
	router.POST("/api/register", userController.Register)
	router.POST("/api/login", userController.Login)

	//task router
	router.POST("/api/tasks", middleware.AuthMiddleware(taskController.Create))
	router.PUT("/api/tasks/:taskId", middleware.AuthMiddleware(taskController.Update))
	router.PATCH("/api/tasks/:taskId", middleware.AuthMiddleware(taskController.Update))
	router.DELETE("/api/tasks/:taskId", middleware.AuthMiddleware(taskController.Delete))
	router.GET("/api/task/:taskId", middleware.AuthMiddleware(taskController.FindById))
	router.GET("/api/tasks/", middleware.AuthMiddleware(taskController.FindAll))

	return router
}
