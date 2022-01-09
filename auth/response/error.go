package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorPayload struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	StatusCode int          `json:"statusCode"`
	Payload    ErrorPayload `json:"payload"`
}

func CreateHttpError(rw http.ResponseWriter, code int, message string) {
	resp := ErrorResponse{
		StatusCode: code,
		Payload:    ErrorPayload{Message: message},
	}

	rw.WriteHeader(code)

	encoder := json.NewEncoder(rw)
	err := encoder.Encode(resp)
	if err != nil {
		fmt.Println("Error encoding response")
	}
}
