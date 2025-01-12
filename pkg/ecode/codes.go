package ecode

var (
	OK                 = New(0, "OK")
	NoLogin            = New(101, "NOT_LOGIN")
	RequestErr         = New(400, "INVALID_ARGUMENT")
	Unauthorized       = New(401, "UNAUTHENTICATED")
	AccessDenied       = New(403, "PERMISSION_DENIED")
	NotFound           = New(404, "NOT_FOUND")
	MethodNotAllowed   = New(405, "METHOD_NOT_ALLOWED")
	Canceled           = New(498, "CANCELED")
	ServerErr          = New(500, "INTERNAL_ERROR")
	ServiceUnavailable = New(503, "UNAVAILABLE")
	Deadline           = New(504, "DEADLINE_EXCEEDED")
	LimitExceed        = New(509, "RESOURCE_EXHAUSTED")

	//用户服务
	VerificationMaxLimit   = New(10001, "已达今日请求上限")
	VerificationCodeFailed = New(10002, "验证码错误")
	MobileHasRegistered    = New(10003, "手机号已注册")
	UserNotExisted         = New(10004, "用户不存在")
	//文章服务
	PutBucketObjectErr = New(20001, "文件上传失败")
)
