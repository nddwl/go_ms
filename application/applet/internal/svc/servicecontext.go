package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"zhihu/application/applet/internal/config"
	"zhihu/application/user/user"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
	UserRpc  user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		BizRedis: redis.MustNewRedis(c.BizRedis),
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
