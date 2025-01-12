package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"
	"zhihu/application/like/rpc/internal/types"

	"zhihu/application/like/rpc/internal/svc"
	"zhihu/application/like/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *pb.ThumbupRequest) (*pb.ThumbupResponse, error) {
	// todo: add your logic here and delete this line
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			logx.Errorf("Thmbup Marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusher.Push(context.Background(), string(data))
		if err != nil {
			logx.Errorf("KqPusher.Push error: %v", err)
		}
	})

	return &pb.ThumbupResponse{}, nil
}
