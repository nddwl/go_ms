package ecode

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strconv"
)

func New(code int, message string) Codes {
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

func JsonCtx(ctx context.Context, w http.ResponseWriter, obj interface{}, err error) {
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
		Resource: obj,
	}
	httpx.WriteJsonCtx(ctx, w, resp.Code, &resp)
}
