package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"zhihu/application/applet/internal/config"
	"zhihu/application/user/user"
	"zhihu/pkg/validator"
)

type ServiceContext struct {
	Config    config.Config
	Validator *validator.Validator
	BizRedis  *redis.Redis
	UserRpc   user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Validator: validator.NewValidator(),
		BizRedis:  redis.MustNewRedis(c.BizRedis),
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
