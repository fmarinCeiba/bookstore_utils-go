package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

func (e restErr) Message() string {
	return e.message
}
func (e restErr) Status() int {
	return e.status
}
func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: [ %v ]", e.message, e.status, e.error, e.causes)
}
func (e restErr) Causes() []interface{} {
	return e.causes
}

type restErr struct {
	message string        `json:"message"`
	status  int           `json:"status"`
	error   string        `json:"error"`
	causes  []interface{} `json:"causes"`
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		message: message,
		status:  status,
		error:   err,
		causes:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var rErr restErr
	if err := json.Unmarshal(bytes, &rErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return rErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusBadRequest,
		error:   "bad_request",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusNotFound,
		error:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusUnauthorized,
		error:   "unauthorized",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	res := restErr{
		message: message,
		status:  http.StatusInternalServerError,
		error:   "internal_server_error",
	}
	if err != nil {
		res.causes = append(res.causes, err.Error())
	}
	return res
}
