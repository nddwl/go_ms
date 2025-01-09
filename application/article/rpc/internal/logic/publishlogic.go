package logic

import (
	"context"
	"time"
	"zhihu/application/article/rpc/internal/model"

	"zhihu/application/article/rpc/internal/svc"
	"zhihu/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	article := model.Article{
		Title:       in.Title,
		Content:     in.Content,
		Cover:       in.Cover,
		Description: in.Description,
		AuthorId:    in.UserId,
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	result, err := l.svcCtx.ArticleModel.Insert(l.ctx, &article)
	if err != nil {
		logx.Errorf("Insert article: %v error: %v", article, err)
		return nil, err
	}
	articleId, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId error: %v", err)
		return nil, err
	}
	return &pb.PublishResponse{ArticleId: articleId}, nil
}
