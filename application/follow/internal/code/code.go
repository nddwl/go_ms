package code

import "zhihu/pkg/ecode"

var (
	FollowedUserIdEmpty = ecode.New(60002, "被关注用户id为空")
	CannotFollowSelf    = ecode.New(60003, "不能关注自己")
	UserIdEmpty         = ecode.New(60004, "用户id为空")
)
