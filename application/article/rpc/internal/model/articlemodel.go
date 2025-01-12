package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zhihu/application/article/rpc/internal/types"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		ArticlesByUserId(ctx context.Context, userId int64, status int32, sortField string,
			likeNum int64, publishTime string) ([]*Article, error)
		UpdateArticleStatus(ctx context.Context, articleId int64, status int64) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}

func (m *customArticleModel) ArticlesByUserId(ctx context.Context, userId int64, status int32, sortField string,
	likeNum int64, publishTime string) ([]*Article, error) {
	var (
		err      error
		query    string
		field    any
		articles []*Article
	)
	if sortField == "like_num" {
		field = likeNum
		query = fmt.Sprintf("select %s from %s where author_id = ? and status = ? and like_num < ? order by %s desc , publish_time desc limit ?", articleRows, m.table, sortField)
	} else {
		field = publishTime
		query = fmt.Sprintf("select %s from %s where author_id = ? and status = ? and publish_time < ? order by %s desc, like_num desc limit ?", articleRows, m.table, sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &articles, query, userId, status, field, types.MaxLimit)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (m *customArticleModel) UpdateArticleStatus(ctx context.Context, articleId int64, status int64) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set status = ? where id = ?", m.table)
		return conn.ExecCtx(ctx, query, status, articleId)
	})
	return err
}
