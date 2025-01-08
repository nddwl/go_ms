package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zhihu/application/user/internal/svc"
	"zhihu/application/user/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &service.FindByIdResponse{UserId: -1}, nil
		} else {
			logx.Errorf("FindOne id: %v error: %v", in.UserId, err)
			return nil, err
		}
	}
	return &service.FindByIdResponse{
		UserId:   user.Id,
		Username: user.Username,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
	}, nil
}
