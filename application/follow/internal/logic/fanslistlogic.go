package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"time"
	"zhihu/application/follow/internal/code"
	"zhihu/application/follow/internal/model"
	"zhihu/application/follow/internal/types"

	"zhihu/application/follow/internal/svc"
	"zhihu/application/follow/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FansListLogic) FansList(in *pb.FansListRequest) (*pb.FansListResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize <= 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor <= 0 {
		in.Cursor = time.Now().Unix()
	}

	var (
		err                     error
		isCache, isEnd          bool
		cursor, fanUserId       int64
		fansUserIds, createTime []int64
		follow                  []*model.Follow
		followCount             []*model.FollowCount
		curPage                 []*pb.FansItem
	)

	fansUserIds, createTime, _ = l.cacheFansList(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(fansUserIds) > 0 {
		isCache = true
		if fansUserIds[len(fansUserIds)-1] == -1 {
			isEnd = true
			fansUserIds = fansUserIds[:len(fansUserIds)-1]
		}
		followCount, err = l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, fansUserIds)
		if err != nil {
			logx.Errorf("FollowCountModel->FindByUserIds in: %v error: %v", in, err)
			return nil, err
		}
	} else {
		follow, err = l.svcCtx.FollowModel.FindByFollowedUserId(l.ctx, in.UserId, types.CacheMaxFansCount)
		if err != nil {
			logx.Errorf("FollowModel->FindByUserId in: %v error: %v", in, err)
			return nil, err
		}
		if len(follow) == 0 {
			return &pb.FansListResponse{}, nil
		}

		var firstFollow []*model.Follow

		if len(follow) < int(in.PageSize) {
			isEnd = true
			firstFollow = follow
		} else {
			firstFollow = follow[:in.PageSize]
		}

		fansUserIds = make([]int64, 0, len(firstFollow))
		createTime = make([]int64, 0, len(firstFollow))
		for _, v := range firstFollow {
			fansUserIds = append(fansUserIds, v.UserID)
			createTime = append(createTime, v.CreateTime.Unix())
		}

		followCount, err = l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, fansUserIds)
		if err != nil {
			logx.Errorf("FollowCountModel->FindByUserIds in: %v error: %v", in, err)
			return nil, err
		}
	}

	fcMap := make(map[int64]*model.FollowCount)
	for i := 0; i < len(followCount); i++ {
		fcMap[followCount[i].UserID] = followCount[i]
	}

	curPage = make([]*pb.FansItem, 0, len(fansUserIds))
	for k, v := range fansUserIds {
		fc := fcMap[v]
		if fc == nil {
			continue
		}
		curPage = append(curPage, &pb.FansItem{
			FansUserId:  v,
			FollowCount: int64(fc.FollowCount),
			FansCount:   int64(fc.FansCount),
			CreateTime:  createTime[k],
		})
	}

	for i := 0; i < len(curPage); i++ {
		if curPage[i].FansUserId == in.FansUserId && curPage[i].CreateTime == in.Cursor {
			curPage = curPage[i+1:]
		}
	}

	if len(curPage) > 0 {
		last := curPage[len(curPage)-1]
		fanUserId = last.FansUserId
		cursor = last.CreateTime
		if cursor < 0 {
			cursor = 0
		}
	}
	resp := &pb.FansListResponse{
		Items:      curPage,
		UserId:     in.UserId,
		Cursor:     cursor,
		IsEnd:      isEnd,
		PageSize:   in.PageSize,
		FansUserId: fanUserId,
	}
	if !isCache {
		threading.GoSafe(func() {
			if len(follow) < types.CacheMaxFansCount && len(follow) > 0 {
				follow = append(follow, &model.Follow{
					UserID:         -1,
					FollowedUserID: -1,
					CreateTime:     time.Unix(0, 0),
				})
			}
			err = l.addFansList(context.Background(), in.UserId, follow)
		})
	}
	return resp, nil
}

func fansKey(userId int64) string {
	return fmt.Sprintf("%s:fans:%d", types.Prefix, userId)
}

func (l *FansListLogic) addFansList(ctx context.Context, userId int64, fans []*model.Follow) error {
	if len(fans) == 0 {
		return nil
	}
	key := fansKey(userId)
	for _, v := range fans {
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, v.CreateTime.Unix(), strconv.FormatInt(v.UserID, 10))
		if err != nil {
			logx.Errorf("ZaddCtx key: %s error: %v", key, err)
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.UserFollowExpireTime)
}

func (l *FansListLogic) cacheFansList(ctx context.Context, userId, cursor int64, pageSize int64) ([]int64, []int64, error) {
	key := fansKey(userId)
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimitCtx key: %s error: %v", key, err)
		return nil, nil, err
	}
	var ids = make([]int64, 0, len(pairs))
	var createTime = make([]int64, 0, len(pairs))
	for _, v := range pairs {
		id, err := strconv.ParseInt(v.Key, 10, 64)
		if err != nil {
			logx.Errorf("ParseInt s: %s error: %v", v.Key, err)
			return nil, nil, err
		}
		ids = append(ids, id)
		createTime = append(createTime, v.Score)
	}
	return ids, createTime, err
}
