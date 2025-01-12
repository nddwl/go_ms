package types

const (
	Prefix            = "article"
	ArticleExpireTime = 3600 * 24 * 3
	MaxLimit          = 200
)

const (
	SortPublishTime = iota
	SortLikeNum
)

const (
	// ArticleStatusPending 待审核
	ArticleStatusPending = iota
	// ArticleStatusNotPass 审核不通过
	ArticleStatusNotPass
	// ArticleStatusVisible 可见
	ArticleStatusVisible
	// ArticleStatusUserDelete 用户删除
	ArticleStatusUserDelete
)
