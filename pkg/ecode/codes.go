package ecode

var (
	Ok         = New(200, "成功")
	ServerErr  = New(500, "服务繁忙")
	BadRequest = New(400, "请求错误")

	VerificationMaxLimit   = New(429, "已达今日请求上限")
	VerificationCodeFailed = New(400, "验证码错误")
	MobileHasRegistered    = New(400, "手机号已注册")
	UserNotExisted         = New(400, "用户不存在")
)
