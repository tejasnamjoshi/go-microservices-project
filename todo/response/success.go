package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	StatusCode int         `json:"statusCode"`
	Payload    interface{} `json:"payload"`
}

func (r *response) CreateSuccessResponse(rw http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{
		StatusCode: http.StatusOK,
		Payload:    data,
	}
	rw.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(rw)
	err := encoder.Encode(resp)
	if err != nil {
		r.CreateHttpError(rw, http.StatusInternalServerError, "Error while encoding response")
		return
	}
}
