package logic

import (
	"context"
	"time"
	"zhihu/application/user/internal/code"
	"zhihu/application/user/internal/model"

	"zhihu/application/user/internal/svc"
	"zhihu/application/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if in.Username == "" {
		return nil, code.RegisterNameEmpty
	}
	user := model.User{
		Username:   in.Username,
		Mobile:     in.Mobile,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &user)
	if err != nil {
		logx.Errorf("Insert user: %v error: %v", user, err)
		return nil, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId error: %v", err)
		return nil, err
	}
	return &pb.RegisterResponse{UserId: userId}, nil
}
