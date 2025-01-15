package logic

import (
	"context"

	"zhihu/application/follow/internal/svc"
	"zhihu/application/follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.FollowResponse{}, nil
}
