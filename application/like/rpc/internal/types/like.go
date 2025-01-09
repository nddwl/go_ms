package types

type (
	ThumbupMsg struct {
		BizId    string `json:"biz_id"`
		ObjId    int64  `json:"obj_id"`
		UserId   int64  `json:"user_id"`
		LikeType int32  `json:"like_type"`
	}
)
