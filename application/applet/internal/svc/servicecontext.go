package svc

import (
	"applet/internal/config"
	"applet/internal/validator"
	"applet/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Validator *validator.Validator
	Redis     *redis.Redis
	UserRpc   service.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Validator: validator.NewValidator(),
		Redis:     redis.MustNewRedis(c.Redis),
		UserRpc:   service.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
	}
}
