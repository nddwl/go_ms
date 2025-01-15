package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"strconv"
	"time"
	"zhihu/application/article/mq/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"zhihu/application/article/mq/internal/svc"
)

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLogic) Consume(ctx context.Context, _, val string) error {
	logx.Infof("Consume msg val: %s", val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume msg val: %s error: %v", val, err)
		return err
	}
	return l.article(ctx, msg)
}

func articleKey(userId string, sortType int32) string {
	return fmt.Sprintf("%s:%s:%d", types.Prefix, userId, sortType)
}

func articleScore(sortType int32, likeNum, publishTime int64) float64 {
	if sortType == types.SortPublishTime {
		return float64(publishTime)
	} else {
		return float64(likeNum) + float64(publishTime)/1e10
	}
}

func (l *ArticleLogic) article(ctx context.Context, msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}
	for _, data := range msg.Data {
		status, _ := strconv.Atoi(data.Status)
		likeKey := articleKey(data.AuthorId, types.SortLikeNum)
		publishKey := articleKey(data.AuthorId, types.SortPublishTime)
		likeNum, _ := strconv.ParseInt(data.LikeNum, 10, 64)
		publishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)
		likeScore := articleScore(types.SortLikeNum, likeNum, publishTime.Unix())
		publishScore := articleScore(types.SortPublishTime, likeNum, publishTime.Unix())
		switch status {
		case types.ArticleStatusVisible:
			b, _ := l.svcCtx.BizRedis.ExistsCtx(ctx, likeKey)
			if b {
				_, err := l.svcCtx.BizRedis.ZaddFloatCtx(ctx, likeKey, likeScore, data.AuthorId)
				if err != nil {
					logx.Errorf("ZaddFloatCtx key: %s value: %s error: %v", likeKey, data.AuthorId, err)
				}
			}
			b, _ = l.svcCtx.BizRedis.ExistsCtx(ctx, publishKey)
			if b {
				_, err := l.svcCtx.BizRedis.ZaddFloatCtx(ctx, publishKey, publishScore, data.AuthorId)
				if err != nil {
					logx.Errorf("ZaddFloatCtx key: %s value: %s error: %v", publishKey, data.AuthorId, err)
				}
			}
		case types.ArticleStatusUserDelete:
			_, err := l.svcCtx.BizRedis.ZremCtx(ctx, likeKey, data.AuthorId)
			if err != nil {
				logx.Errorf("ZremCtx key: %s value: %s", likeKey, data.AuthorId)
			}
			_, err = l.svcCtx.BizRedis.ZremCtx(ctx, publishKey, data.AuthorId)
			if err != nil {
				logx.Errorf("ZremCtx key: %s value: %s", publishKey, data.AuthorId)
			}
		}

	}
	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumer, NewArticleLogic(ctx, svcCtx)),
	}
}
