package response

import (
	"encoding/json"
	"net/http"
)

type ErrorPayload struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	StatusCode int          `json:"statusCode"`
	Payload    ErrorPayload `json:"payload"`
}

func (r *response) SendErrorResponse(rw http.ResponseWriter, code int, message string) {
	resp := ErrorResponse{
		StatusCode: code,
		Payload:    ErrorPayload{Message: message},
	}

	rw.WriteHeader(code)

	encoder := json.NewEncoder(rw)
	err := encoder.Encode(resp)
	if err != nil {
		r.Logger.Error(err.Error())
	}
}
