package interceptor

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"zhihu/pkg/ecode"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			sts := status.Convert(err)
			codes := ecode.ToCodes(sts)
			err = errors.WithMessage(codes, codes.Message())
		}
		return err
	}
}
