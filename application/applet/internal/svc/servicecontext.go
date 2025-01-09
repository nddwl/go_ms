package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"zhihu/application/applet/internal/config"
	"zhihu/application/applet/service"
	"zhihu/pkg/validator"
)

type ServiceContext struct {
	Config    config.Config
	Validator *validator.Validator
	Redis     *redis.Redis
	UserRpc   service.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	//启用自定义错误处理
	//zrpc.WithUnaryClientInterceptor(ecode.ClientErrorInterceptor())
	userRpc := service.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn())
	return &ServiceContext{
		Config:    c,
		Validator: validator.NewValidator(),
		Redis:     redis.MustNewRedis(c.Redis),
		UserRpc:   userRpc,
	}
}
