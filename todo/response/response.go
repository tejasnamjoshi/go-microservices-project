package response

import (
	"go-todo/auth/logging"
	"net/http"
)

type Response interface {
	CreateSuccessResponse(rw http.ResponseWriter, data interface{})
	CreateHttpError(rw http.ResponseWriter, code int, message string)
}

type response struct {
	Logger logging.Logger
}

func NewResponse(logger logging.Logger) *response {
	return &response{Logger: logger}
}
