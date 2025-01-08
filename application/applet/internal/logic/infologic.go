package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"zhihu/application/applet/internal/svc"
	"zhihu/application/applet/internal/types"
	"zhihu/application/applet/service"
	"zhihu/pkg/ecode"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.InfoRequest) (resp *types.InfoResponse, err error) {
	if err = l.svcCtx.Validator.Struct(req); err != nil {
		err = ecode.BadRequest
		return
	}
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		logx.Errorf("ctx.Value(%s) error: %v", types.UserIdKey, err)
		return nil, err
	}
	if userId == 0 {
		return &types.InfoResponse{}, nil
	}
	IdResp, err := l.svcCtx.UserRpc.FindById(l.ctx, &service.FindByIdRequest{UserId: userId})
	if err != nil {
		logx.Errorf("userRpc->FindById userId: %d error: %v", userId, err)
		return nil, err
	}
	return &types.InfoResponse{
		UserId:   IdResp.UserId,
		Username: IdResp.Username,
		Avatar:   IdResp.Avatar,
	}, nil
}
