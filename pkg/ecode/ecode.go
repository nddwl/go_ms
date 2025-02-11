package ecode

import (
	"context"
	"errors"
	"strconv"
)

var _codes = make(map[int]struct{})

func New(code int, message string) Code {
	if _, ok := _codes[code]; ok {
		panic("code重复")
	}
	return Code{
		code:    code,
		message: message,
	}
}

type Response struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Resource any    `json:"resource"`
}

type Codes interface {
	error
	Code() int
	Message() string
	Details() []interface{}
}

var _ Codes = Code{}

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

func (c Code) Details() []interface{} {
	return nil
}

func Cause(err error) Codes {
	var e Codes
	switch {
	case err == nil:
		e = OK
	case errors.As(err, &e):
	case errors.Is(err, context.Canceled):
		e = Canceled
	case errors.Is(err, context.DeadlineExceeded):
		e = Deadline
	default:
		e = ServerErr
	}
	return e
}

func ErrorHandler() func(err error) (int, any) {
	return func(err error) (int, any) {
		e := Cause(err)
		resp := Response{
			Code:     e.Code(),
			Message:  e.Message(),
			Resource: nil,
		}
		return 200, &resp
	}
}

func OkHandler() func(ctx context.Context, a any) any {
	return func(ctx context.Context, a any) any {
		return &Response{
			Code:     OK.Code(),
			Message:  OK.Message(),
			Resource: a,
		}
	}
}
