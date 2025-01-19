package logic

import (
	"context"
	"errors"
	"zhihu/application/article/rpc/internal/code"
	"zhihu/application/article/rpc/internal/model"
	"zhihu/application/article/rpc/internal/svc"
	"zhihu/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailRequest) (*pb.ArticleDetailResponse, error) {
	if in.ArticleId <= 0 {
		return nil, code.ArticleIdInvalid
	}
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return &pb.ArticleDetailResponse{Article: &pb.ArticleItem{Id: -1}}, nil
		}
		return nil, err
	}
	return &pb.ArticleDetailResponse{
		Article: &pb.ArticleItem{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			Cover:       article.Cover,
			Description: article.Description,
			CommentNum:  article.CommentNum,
			LikeNum:     article.LikeNum,
			CollectNum:  article.CollectNum,
			ViewNum:     article.ViewNum,
			ShareNum:    article.ShareNum,
			TagIds:      article.TagIds,
			PublishTime: article.PublishTime.Unix(),
			UpdateTime:  article.UpdateTime.Unix(),
			AuthorId:    article.AuthorId,
		},
	}, nil
}
