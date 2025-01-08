package ecode

var (
	Ok         = New(200, "成功")
	ServerErr  = New(500, "服务繁忙")
	BadRequest = New(400, "请求错误")

	VerificationMaxLimit   = New(100100, "已达今日请求上限")
	VerificationCodeFailed = New(100101, "验证码错误")
	MobileHasRegistered    = New(100200, "手机号已注册")
	UserNotExisted         = New(100201, "用户不存在")
)
