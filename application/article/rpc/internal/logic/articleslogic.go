package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"slices"
	"strconv"
	"time"
	"zhihu/application/article/rpc/internal/code"
	"zhihu/application/article/rpc/internal/model"
	"zhihu/application/article/rpc/internal/svc"
	"zhihu/application/article/rpc/internal/types"
	"zhihu/application/article/rpc/pb"
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.SortType != types.SortLikeNum && in.SortType != types.SortPublishTime {
		return nil, code.SortTypeInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor <= 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = float64(time.Now().Unix())
		} else {
			in.Cursor = float64(types.DefaultSortLikeCursor)
		}
	}
	var (
		err                    error
		isCache, isEnd         bool
		likeNum                int64
		sortField, publishTime string
		curPage                []*pb.ArticleItem
		articles               []*model.Article
	)
	if in.SortType == types.SortPublishTime {
		sortField = "publish_time"
		publishTime = time.Unix(int64(in.Cursor), 0).Format("2006-01-02 15:04:05")
	} else {
		sortField = "like_num"
		likeNum = int64(in.Cursor)
	}
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.PageSize, in.Cursor, in.SortType)
	if len(articleIds) > 0 {
		isCache = true
		if articleIds[len(articleIds)-1] == -1 {
			articleIds = articleIds[:len(articleIds)-1]
			isEnd = true
		}
		articles, err = l.articlesByUserId(l.ctx, articleIds)
		if err != nil {
			logx.Errorf("articlesById error: %v", err)
			return nil, err
		}
		sortArticles(articles, in.SortType)
		curPage = make([]*pb.ArticleItem, 0, len(articles))
		for i := 0; i < len(articles); i++ {
			article := articles[i]
			curPage = append(curPage, &pb.ArticleItem{
				Id:          article.Id,
				Title:       article.Title,
				Content:     article.Content,
				Cover:       article.Cover,
				Description: article.Description,
				CommentNum:  article.CommentNum,
				LikeNum:     article.LikeNum,
				CollectNum:  article.CollectNum,
				ViewNum:     article.ViewNum,
				ShareNum:    article.ShareNum,
				TagIds:      article.TagIds,
				PublishTime: article.PublishTime.Unix(),
				UpdateTime:  article.UpdateTime.Unix(),
				AuthorId:    article.AuthorId,
			})
		}
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(articleKey(in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.ArticlesByUserId(l.ctx, in.UserId, types.ArticleStatusVisible, sortField, likeNum, publishTime)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d softType:%d error: %v", in.UserId, in.SortType, err)
			return nil, err
		}
		if v == nil {
			return &pb.ArticlesResponse{}, nil
		}
		articles = v.([]*model.Article)
		var tempArticles []*model.Article
		if len(articles) > int(in.PageSize) {
			tempArticles = articles[:in.PageSize]
		} else {
			tempArticles = articles
			isEnd = true
		}
		for i := 0; i < len(tempArticles); i++ {
			article := articles[i]
			curPage = append(curPage, &pb.ArticleItem{
				Id:          article.Id,
				Title:       article.Title,
				Content:     article.Content,
				Cover:       article.Cover,
				Description: article.Description,
				CommentNum:  article.CommentNum,
				LikeNum:     article.LikeNum,
				CollectNum:  article.CollectNum,
				ViewNum:     article.ViewNum,
				ShareNum:    article.ShareNum,
				TagIds:      article.TagIds,
				PublishTime: article.PublishTime.Unix(),
				UpdateTime:  article.UpdateTime.Unix(),
				AuthorId:    article.AuthorId,
			})
		}
	}

	var (
		cursor    = in.Cursor
		articleId = in.ArticleId
	)
	for i := 0; i < len(curPage); i++ {
		article := curPage[i]
		if in.SortType == types.SortPublishTime {
			if article.Id == in.ArticleId && article.PublishTime == int64(in.Cursor) {
				curPage = curPage[i+1:]
				break
			}
		} else {
			if article.Id == in.ArticleId && article.LikeNum == int64(in.Cursor) {
				curPage = curPage[i+1:]
				break
			}
		}
	}
	if len(curPage) > 0 {
		lastArticle := curPage[len(curPage)-1]
		cursor = articleScore(in.SortType, lastArticle.LikeNum, lastArticle.PublishTime)
		articleId = lastArticle.Id
	}

	resp := &pb.ArticlesResponse{
		Articles:  curPage,
		IsEnd:     isEnd,
		UserId:    in.UserId,
		PageSize:  in.PageSize,
		SortType:  in.SortType,
		Cursor:    cursor,
		ArticleId: articleId,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(articles) < types.MaxLimit && len(articles) > 0 {
				articles = append(articles, &model.Article{Id: -1, PublishTime: time.Unix(0, 0)})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheArticles error: %v", err)
			}
		})
	}
	return resp, nil
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, articles []*model.Article, userId int64, sortType int32) error {
	if len(articles) < 0 {
		return nil
	}
	key := articleKey(userId, sortType)
	for i := 0; i < len(articles); i++ {
		_, err := l.svcCtx.BizRedis.ZaddFloatCtx(ctx, key, articleScore(sortType, articles[i].LikeNum, articles[i].PublishTime.Unix()), strconv.FormatInt(articles[i].Id, 10))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.ArticleExpireTime)
}

func articleKey(userId int64, sortType int32) string {
	return fmt.Sprintf("%s:%d:%d", types.Prefix, userId, sortType)
}

func articleScore(sortType int32, likeNum, publishTime int64) float64 {
	if sortType == types.SortPublishTime {
		return float64(publishTime)
	} else {
		return float64(likeNum) + float64(publishTime)/1e10
	}
}

func sortArticles(articles []*model.Article, sortType int32) {
	var cmpFun func(a, b *model.Article) int
	if sortType == types.SortPublishTime {
		cmpFun = func(a, b *model.Article) int {
			if a.PublishTime.Unix() < b.PublishTime.Unix() {
				return 1
			}
			if a.PublishTime.Unix() == b.PublishTime.Unix() {
				if a.LikeNum < b.LikeNum {
					return 1
				}
				return -1
			}
			return -1
		}
	} else {
		cmpFun = func(a, b *model.Article) int {
			if a.LikeNum < b.LikeNum {
				return 1
			}
			if a.LikeNum == b.LikeNum {
				if a.PublishTime.Unix() < b.PublishTime.Unix() {
					return 1
				}
				return -1
			}
			return -1
		}
	}
	slices.SortFunc(articles, cmpFun)
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, userId, pageSize int64, cursor float64, sortType int32) ([]int64, error) {
	var (
		err        error
		key        string
		pairs      []redis.FloatPair
		articleIds []int64
	)
	key = articleKey(userId, sortType)
	ok, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("Exists key: %s error: %v", key, err)
	}
	if ok {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, types.ArticleExpireTime)
		if err != nil {
			logx.Errorf("Expire key: %s error: %v", key, err)
		}
	}
	pairs, err = l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresByFloatAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	articleIds = make([]int64, 0, len(pairs))
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		articleIds = append(articleIds, id)
	}
	return articleIds, nil
}

func (l *ArticlesLogic) articlesByUserId(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](func(source chan<- int64) {
		for _, id := range articleIds {
			if id == -1 {
				continue
			}
			source <- id
		}
	}, func(articleId int64, writer mr.Writer[*model.Article], cancel func(error)) {
		article, err := l.svcCtx.ArticleModel.FindOne(ctx, articleId)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(article)
	}, func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
		var articles []*model.Article
		for article := range pipe {
			articles = append(articles, article)
		}
		writer.Write(articles)
	})
	if err != nil {
		return nil, err
	}
	return articles, err
}
