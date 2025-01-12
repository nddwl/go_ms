package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"zhihu/pkg/ecode"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		err = ecode.FromError(err).Err()
		return resp, err
	}
}
