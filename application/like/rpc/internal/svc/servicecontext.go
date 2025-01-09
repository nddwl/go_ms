package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"zhihu/application/like/rpc/internal/config"
)

type ServiceContext struct {
	Config   config.Config
	KqPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		KqPusher: kq.NewPusher(c.KqPusher.Brokers, c.KqPusher.Topic),
	}
}
