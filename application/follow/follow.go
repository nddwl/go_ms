package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"zhihu/pkg/ecode"
	"zhihu/pkg/interceptor"

	"zhihu/application/follow/internal/config"
	"zhihu/application/follow/internal/server"
	"zhihu/application/follow/internal/svc"
	"zhihu/application/follow/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/follow.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFollowServer(grpcServer, server.NewFollowServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	httpx.SetErrorHandler(ecode.ErrorHandler())
	httpx.SetOkHandler(ecode.OkHandler())

	s.AddUnaryInterceptors(interceptor.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
