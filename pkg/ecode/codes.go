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
)
