package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"zhihu/pkg/orm"
)

type Config struct {
	zrpc.RpcServerConf
	DB       orm.Config
	BizRedis redis.RedisConf
}
