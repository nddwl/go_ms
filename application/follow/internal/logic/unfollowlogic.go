package logic

import (
	"context"
	"gorm.io/gorm"
	"zhihu/application/follow/internal/code"
	"zhihu/application/follow/internal/model"
	"zhihu/application/follow/internal/types"

	"zhihu/application/follow/internal/svc"
	"zhihu/application/follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFollowLogic) UnFollow(in *pb.UnFollowRequest) (*pb.UnFollowResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdEmpty
	}
	if in.FollowedUserId <= 0 {
		return nil, code.FollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId {
		return nil, code.CannotFollowSelf
	}

	follow, err := l.svcCtx.FollowModel.FindByUserIdAndFollowedUserId(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		logx.Errorf("FollowModel->FindByUserIdAndFollowedUserId in: %v error: %v", in, err)
		return nil, err
	}
	if follow == nil || follow.FollowStatus == types.FollowStatusUnfollow {
		return &pb.UnFollowResponse{}, nil
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{
			"follow_status": types.FollowStatusUnfollow,
		})
		if err != nil {
			return err
		}
		err = model.NewFollowCountModel(tx).DecrFollowCount(l.ctx, follow.UserID)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).DecrFansCount(l.ctx, follow.FollowedUserID)
	})
	if err != nil {
		logx.Errorf("Transaction error: %v", err)
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, fansKey(in.FollowedUserId), in.UserId)
	if err != nil {
		logx.Errorf("ZremCtx error: %v", err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, followedKey(in.UserId), in.FollowedUserId)
	if err != nil {
		logx.Errorf("ZremCtx error: %v", err)
		return nil, err
	}
	return &pb.UnFollowResponse{}, nil
}
