package response

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = 200
	StatusError = 400
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidateErrors(errors validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "tagName":
			errMsg = append(errMsg, fmt.Sprintf("filed %s is a name failed"))
		default:
			errMsg = append(errMsg, fmt.Sprintf("%s is not valid"))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
