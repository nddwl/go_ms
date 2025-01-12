package ecode

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
	"strconv"
	"zhihu/pkg/ecode/types"
)

var _ Codes = &Status{}

type Status struct {
	sts *types.Status
}

func (s *Status) Error() string {
	return fmt.Sprintf("Error%d:%s", s.sts.Code, s.sts.Message)
}

func (s *Status) Code() int {
	return int(s.sts.Code)
}

func (s *Status) Message() string {
	return s.sts.Message
}

func (s *Status) Details() []any {
	if s == nil || s.sts == nil {
		return nil
	}
	details := make([]any, 0, len(s.sts.Details))
	for _, d := range s.sts.Details {
		detail, err := d.UnmarshalNew()
		if err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, protoadapt.MessageV1Of(detail))
	}
	return details
}

func FromCode(c Code) *Status {
	return &Status{sts: &types.Status{
		Code:    int32(c.code),
		Message: c.message,
		Details: nil,
	}}
}

func FromProto(sts *types.Status) Codes {
	return &Status{sts: sts}
}

func FromError(err error) *status.Status {
	var c Codes
	if errors.As(err, &c) {
		sts, e := FromCodes(c)
		if e == nil {
			return sts
		}
	}
	var sts *status.Status
	switch {
	case errors.Is(err, context.Canceled):
		sts, _ = FromCodes(Canceled)
	case errors.Is(err, context.DeadlineExceeded):
		sts, _ = FromCodes(Deadline)
	default:
		sts, _ = status.FromError(err)
	}
	return sts
}

func FromCodes(c Codes) (*status.Status, error) {
	var s *Status
	var code Code
	switch {
	case errors.As(c, &s):
	case errors.As(c, &code):
		s = FromCode(code)
	default:
		s = &Status{sts: &types.Status{
			Code:    int32(c.Code()),
			Message: c.Message(),
		}}
	}
	return status.New(codes.Unknown, strconv.Itoa(s.Code())).WithDetails(s.sts)
}

func ToCodes(s *status.Status) Codes {
	details := s.Details()
	for _, v := range details {
		if c, ok := v.(*types.Status); ok {
			return FromProto(c)
		}
	}
	return GrpcCodesTOCodes(s)
}

func GrpcCodesTOCodes(s *status.Status) Codes {
	grpcCode := s.Code()
	switch grpcCode {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return RequestErr
	case codes.NotFound:
		return NotFound
	case codes.PermissionDenied:
		return AccessDenied
	case codes.Unauthenticated:
		return Unauthorized
	case codes.ResourceExhausted:
		return LimitExceed
	case codes.Unimplemented:
		return MethodNotAllowed
	case codes.DeadlineExceeded:
		return Deadline
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.Unknown:
		return Code{
			code:    ServerErr.code,
			message: s.Message(),
		}
	}

	return ServerErr
}
