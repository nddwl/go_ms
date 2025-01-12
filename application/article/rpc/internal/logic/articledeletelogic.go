package logic

import (
	"context"
	"zhihu/application/article/rpc/internal/types"

	"zhihu/application/article/rpc/internal/svc"
	"zhihu/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(in *pb.ArticleDeleteRequest) (*pb.ArticleDeleteResponse, error) {
	err := l.svcCtx.ArticleModel.UpdateArticleStatus(l.ctx, in.ArticleId, types.ArticleStatusUserDelete)
	if err != nil {
		logx.Errorf("UpdateArticleStatus req: %v error: %v", in, err)
		return nil, err
	}
	likeKey := articleKey(in.UserId, types.SortLikeNum)
	publishKey := articleKey(in.UserId, types.SortPublishTime)
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeKey, in.ArticleId)
	if err != nil {
		logx.Errorf("ZremCtx key: %s value: %d", likeKey, in.ArticleId)
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, publishKey, in.ArticleId)
	if err != nil {
		logx.Errorf("ZremCtx key: %s value: %d", publishKey, in.ArticleId)
	}
	return &pb.ArticleDeleteResponse{}, nil
}
