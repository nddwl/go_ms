package code

import "zhihu/pkg/ecode"

var (
	SortTypeInvalid         = ecode.New(30001, "排序类型无效")
	UserIdInvalid           = ecode.New(30002, "用户ID无效")
	ArticleTitleCantEmpty   = ecode.New(30003, "文章标题不能为空")
	ArticleContentCantEmpty = ecode.New(30004, "文章内容不能为空")
	ArticleIdInvalid        = ecode.New(30005, "文章ID无效")
)
