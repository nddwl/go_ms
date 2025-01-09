package svc

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/zrpc"
	"zhihu/application/article/api/internal/config"
	"zhihu/application/article/rpc/article"
)

type ServiceContext struct {
	Config     config.Config
	ArticleRpc article.Article
	Oss        *oss.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	oc, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret, oss.Timeout(c.Oss.ConnectTimeout, c.Oss.ReadWriteTimeout))
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		ArticleRpc: article.NewArticle(zrpc.MustNewClient(c.ArticleRpc)),
		Oss:        oc,
	}
}
