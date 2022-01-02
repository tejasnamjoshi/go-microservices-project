package utils

import "net/http"

func CreateHttpError(rw http.ResponseWriter, code int, message string) {
	rw.WriteHeader(code)
	rw.Write([]byte(message))
}
