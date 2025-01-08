package ecode

import (
	"errors"
	"strconv"
)

var _codes = make(map[int]struct{})

func New(code int, message string) Codes {
	if _, ok := _codes[code]; ok {
		panic("code重复")
	}
	return Code{
		code:    code,
		message: message,
	}
}

type Response struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Resource interface{} `json:"resource"`
}

type Codes interface {
	error
	Code() int
	Message() string
}

type Code struct {
	code    int
	message string
}

func (c Code) Code() int {
	return c.code
}

func (c Code) Message() string {
	return c.message
}

func (c Code) Error() string {
	return "Error" + strconv.Itoa(c.code) + ":" + c.message
}

func ErrorHandler() func(err error) (int, any) {
	return func(err error) (int, any) {
		var e Codes
		switch {
		case err == nil:
			e = Ok
		case errors.As(err, &e):
		default:
			e = ServerErr
		}
		resp := Response{
			Code:     e.Code(),
			Message:  e.Message(),
			Resource: nil,
		}
		return 200, &resp
	}
}
