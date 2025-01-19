package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"zhihu/application/follow/internal/config"
	"zhihu/application/follow/internal/model"
	"zhihu/pkg/orm"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	FollowModel      *model.FollowModel
	FollowCountModel *model.FollowCountModel
	BizRedis         *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(c.DB)
	return &ServiceContext{
		Config:           c,
		DB:               db,
		FollowModel:      model.NewFollowModel(db.DB),
		FollowCountModel: model.NewFollowCountModel(db.DB),
		BizRedis:         redis.MustNewRedis(c.BizRedis),
	}
}
