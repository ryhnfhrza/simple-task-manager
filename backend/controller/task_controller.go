package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TaskController interface {
	Create(writer http.ResponseWriter, reqest *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, reqest *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, reqest *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, reqest *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, reqest *http.Request, params httprouter.Params)
}
