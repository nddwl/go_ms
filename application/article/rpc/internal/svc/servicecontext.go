package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zhihu/application/article/rpc/internal/config"
	"zhihu/application/article/rpc/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(sqlx.NewMysql(c.DataSource)),
	}
}
