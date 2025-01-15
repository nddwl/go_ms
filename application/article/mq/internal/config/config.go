package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	KqConsumer        kq.KqConf
	ArticleKqConsumer kq.KqConf
	DataSource        string
	BizRedis          redis.RedisConf
}
