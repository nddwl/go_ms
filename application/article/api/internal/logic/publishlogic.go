package logic

import (
	"context"
	"encoding/json"
	"zhihu/application/article/rpc/pb"
	"zhihu/pkg/ecode"
	"zhihu/pkg/validator"

	"zhihu/application/article/api/internal/svc"
	"zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value(%s) error: %v", types.UserIdKey, err)
		return nil, err
	}
	if err = validator.Struct(req); err != nil {
		return nil, ecode.RequestErr
	}
	publish := pb.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	}
	publishResp, err := l.svcCtx.ArticleRpc.Publish(l.ctx, &publish)
	if err != nil {
		logx.Errorf("ArticleRpc->Publish req: %v error: %v", req, err)
		return nil, err
	}
	return &types.PublishResponse{
		ArticleId: publishResp.ArticleId,
	}, nil
}
