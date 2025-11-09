package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/ryhnfhrza/simple-task-manager/app"
	"github.com/ryhnfhrza/simple-task-manager/controller"
	"github.com/ryhnfhrza/simple-task-manager/exception"
	"github.com/ryhnfhrza/simple-task-manager/helper"
	"github.com/ryhnfhrza/simple-task-manager/repository"
	"github.com/ryhnfhrza/simple-task-manager/service"
	"github.com/ryhnfhrza/simple-task-manager/util"
)

func main() {

	envPath := filepath.Join("..", ".env")

	if p := os.Getenv("CONFIG_PATH"); p != "" {
		envPath = p
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: gagal memuat %s: %v", envPath, err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	DB := app.NewDB()

	validate := validator.New()
	util.RegisterValidations(validate)

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, DB, validate)
	userController := controller.NewUserController(userService)

	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, DB, validate)
	taskController := controller.NewTaskController(taskService)

	router := app.NewRouter(userController, taskController)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    ":" + port,
		Handler: app.CORS(router),
	}

	log.Printf("Server running on port %s", port)
	err := server.ListenAndServe()
	helper.PanicIfError(err)

}
