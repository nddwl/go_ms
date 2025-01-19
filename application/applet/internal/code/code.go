package code

import "zhihu/pkg/ecode"

var (
	VerificationMaxLimit   = ecode.New(10001, "已达今日请求上限")
	VerificationCodeFailed = ecode.New(10002, "验证码错误")
	MobileHasRegistered    = ecode.New(10003, "手机号已注册")
	UserNotExisted         = ecode.New(10004, "用户不存在")
)
