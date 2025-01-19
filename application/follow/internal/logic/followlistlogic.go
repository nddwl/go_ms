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

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
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
		err                         error
		isEnd, isCache              bool
		cursor, followedUserId      int64
		followedUserIds, createTime []int64
		followed                    []*model.Follow
		followedCount               []*model.FollowCount
		curPage                     []*pb.FollowItem
	)

	followedUserIds, createTime, _ = l.cacheFollowList(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followedUserIds) > 0 {
		isCache = true
		if followedUserIds[len(followedUserIds)-1] == -1 {
			isEnd = true
			followedUserIds = followedUserIds[:len(followedUserIds)-1]
		}
		followedCount, err = l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, followedUserIds)
		if err != nil {
			logx.Errorf("FollowCountModel->FindByUserIds in: %v error: %v", in, err)
			return nil, err
		}
	} else {
		followed, err = l.svcCtx.FollowModel.FindByUserId(l.ctx, in.UserId, types.CacheMaxFollowCount)
		if err != nil {
			logx.Errorf("FollowModel->FindByFollowedUserId in: %v error: %v", in, err)
			return nil, err
		}
		if len(followed) == 0 {
			return &pb.FollowListResponse{}, nil
		}

		var firstFollowed []*model.Follow
		if len(followed) < int(in.Cursor) {
			isEnd = true
			firstFollowed = followed
		} else {
			firstFollowed = followed[:in.Cursor]
		}

		followedUserIds = make([]int64, 0, len(firstFollowed))
		createTime = make([]int64, 0, len(firstFollowed))
		for _, v := range firstFollowed {
			followedUserIds = append(followedUserIds, v.FollowedUserID)
			createTime = append(createTime, v.CreateTime.Unix())
		}
		followedCount, err = l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, followedUserIds)
		if err != nil {
			logx.Errorf("FollowCountModel->FindByUserIds in: %v error: %v", in, err)
			return nil, err
		}
	}
	fcMap := make(map[int64]*model.FollowCount)
	for i := 0; i < len(followedCount); i++ {
		fcMap[followedCount[i].UserID] = followedCount[i]
	}
	curPage = make([]*pb.FollowItem, 0, len(followedUserIds))
	for k, v := range followedUserIds {
		fc := fcMap[v]
		if fc == nil {
			continue
		}
		curPage = append(curPage, &pb.FollowItem{
			FollowedUserId: v,
			FollowCount:    int64(fc.FollowCount),
			FansCount:      int64(fc.FansCount),
			CreateTime:     createTime[k],
		})
	}

	for i := 0; i < len(curPage); i++ {
		if in.FollowedUserId == curPage[i].FollowedUserId && in.Cursor == curPage[i].CreateTime {
			curPage = curPage[i+1:]
		}
	}

	if len(curPage) > 0 {
		last := curPage[len(curPage)-1]
		followedUserId = last.FollowedUserId
		cursor = last.CreateTime
		if cursor < 0 {
			cursor = 0
		}
	}

	resp := &pb.FollowListResponse{
		Items:          curPage,
		UserId:         in.UserId,
		Cursor:         cursor,
		IsEnd:          isEnd,
		PageSize:       in.PageSize,
		FollowedUserId: followedUserId,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(followed) < types.CacheMaxFollowCount && len(followed) > 0 {
				followed = append(followed, &model.Follow{
					UserID:         -1,
					FollowedUserID: -1,
					CreateTime:     time.Unix(0, 0),
				})
			}
			err = l.addFollowList(context.Background(), in.UserId, followed)
		})
	}

	return resp, nil
}

func followedKey(userId int64) string {
	return fmt.Sprintf("%s:follow:%d", types.Prefix, userId)
}

func (l *FollowListLogic) cacheFollowList(ctx context.Context, userId, cursor int64, pageSize int64) ([]int64, []int64, error) {
	key := followedKey(userId)
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimitCtx key: %s error: %v", key, err)
		return nil, nil, err
	}
	var ids = make([]int64, 0, len(pairs))
	var createTime = make([]int64, 0, len(pairs))
	for _, v := range pairs {
		createTime = append(createTime, v.Score)
		id, err := strconv.ParseInt(v.Key, 10, 64)
		if err != nil {
			logx.Errorf("ParseInt pair: %v error: %v", v, err)
			return nil, nil, err
		}
		ids = append(ids, id)
	}
	return ids, createTime, nil
}

func (l *FollowListLogic) addFollowList(ctx context.Context, userId int64, followed []*model.Follow) error {
	if len(followed) == 0 {
		return nil
	}
	key := followedKey(userId)
	for _, v := range followed {
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, v.CreateTime.Unix(), strconv.FormatInt(v.FollowedUserID, 10))
		if err != nil {
			logx.Errorf("ZaddCtx key: %s error: %v", key, err)
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.UserFollowExpireTime)
}
