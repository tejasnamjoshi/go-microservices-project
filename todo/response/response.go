package response

import (
	"go-todo/todo/logging"
	"net/http"
)

type Response interface {
	SendSuccessResponse(rw http.ResponseWriter, data interface{})
	SendErrorResponse(rw http.ResponseWriter, code int, message string)
}

type response struct {
	Logger logging.Logger
}

func NewResponse(logger logging.Logger) *response {
	return &response{Logger: logger}
}
