package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"zhihu/application/article/mq/internal/config"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		BizRedis: redis.MustNewRedis(c.BizRedis),
	}
}
