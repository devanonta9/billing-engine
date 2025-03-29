package utils

import (
	"encoding/json"
	"net/http"
	"os"

	constants "billing-engine/domain/constants"
)

type ResponseHTTP struct {
	StatusCode int
	Response   ResponseData
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Detail  any    `json:"data,omitempty"`
}

var (
	ErrRespServiceMaintance = ResponseHTTP{
		StatusCode: http.StatusServiceUnavailable,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespUnauthorize = ResponseHTTP{
		StatusCode: http.StatusUnauthorized,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespAuthInvalid = ResponseHTTP{
		StatusCode: http.StatusUnauthorized,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespBadRequest = ResponseHTTP{
		StatusCode: http.StatusBadRequest,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespInternalServer = ResponseHTTP{
		StatusCode: http.StatusServiceUnavailable,
		Response:   ResponseData{Status: constants.Fail}}
)

func WriteResponse(res http.ResponseWriter, resp ResponseData, code int) {
	res.Header().Set("Content-Type", "application/json")
	r, _ := json.Marshal(resp)

	res.WriteHeader(code)
	res.Write(r)
}

type Error struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

func NewError(id string, status string, title string) *Error {
	return &Error{
		Id:     id,
		Status: status,
		Title:  title,
	}
}

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
