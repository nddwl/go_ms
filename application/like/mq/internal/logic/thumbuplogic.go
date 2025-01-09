package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"

	"github.com/zeromicro/go-zero/core/logx"
	"zhihu/application/like/mq/internal/svc"
)

type ThumbupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThumbupLogic) Consume(ctx context.Context, key string, value string) error {
	// todo: add your logic here and delete this line
	fmt.Println(key, value)
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumer, NewThumbupLogic(ctx, svcCtx)),
	}
}
