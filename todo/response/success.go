package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	StatusCode int         `json:"statusCode"`
	Payload    interface{} `json:"payload"`
}

// Encodes the data onto the response-writer.
// Invokes the errorhandler if the response cannot be encoded.
func (r *response) SendSuccessResponse(rw http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{
		StatusCode: http.StatusOK,
		Payload:    data,
	}
	rw.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(rw)
	err := encoder.Encode(resp)
	if err != nil {
		r.SendErrorResponse(rw, http.StatusInternalServerError, "Error while encoding response")
		return
	}
}
