package apiserver

import "net/http"

func GetErrorCode(err error) int {
	return http.StatusInternalServerError
}
