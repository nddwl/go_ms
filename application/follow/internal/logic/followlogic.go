package logic

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"time"
	"zhihu/application/follow/internal/code"
	"zhihu/application/follow/internal/model"
	"zhihu/application/follow/internal/types"

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
		logx.Errorf("FollowModel->FindByUserIdAndFollowedUserId req: %v error: %v", in, err)
		return nil, err
	}
	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}
	t := time.Now()
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if follow != nil {
			err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{"follow_status": types.FollowStatusFollow})
		} else {
			err = model.NewFollowModel(tx).Insert(l.ctx, &model.Follow{
				UserID:         in.UserId,
				FollowedUserID: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
				CreateTime:     t,
				UpdateTime:     t,
			})
		}
		if err != nil {
			return err
		}
		err = model.NewFollowCountModel(tx).IncrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).IncrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		logx.Errorf("Transaction error: %v", err)
		return nil, err
	}

	fansK := fansKey(in.FollowedUserId)
	followK := followedKey(in.UserId)

	b, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, fansK)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", fansK, err)
		return nil, err
	}
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, fansK, t.Unix(), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			logx.Errorf("ZaddCtx key: %s error: %v", fansK, err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, fansK, 0, -(types.CacheMaxFansCount + 1))
		if err != nil {
			logx.Errorf("ZremrangebyrankCtx key: %s error: %v", fansK, err)
			return nil, err
		}
	}
	b, err = l.svcCtx.BizRedis.ExistsCtx(l.ctx, followK)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", fansK, err)
		return nil, err
	}
	if b {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, followK, t.Unix(), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			logx.Errorf("ZaddCtx key: %s error: %v", fansK, err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, followK, 0, -(types.CacheMaxFollowCount + 1))
		if err != nil {
			logx.Errorf("ZremrangebyrankCtx key: %s error: %v", followK, err)
			return nil, err
		}
	}
	return &pb.FollowResponse{}, nil
}
