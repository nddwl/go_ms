package logic

import (
	"context"
	"net/http"
	"path/filepath"
	"zhihu/pkg/ecode"
	"zhihu/pkg/utils"

	"zhihu/application/article/api/internal/svc"
	"zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	_ = req.ParseMultipartForm(types.MaxFileSize)
	file, header, err := req.FormFile(types.FileName)
	if err != nil {
		return nil, ecode.RequestErr
	}
	bucket, err := l.svcCtx.Oss.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		logx.Errorf("Oss.Bucket error: %v", err)
		return nil, ecode.ServerErr
	}
	key := utils.GenerateUUID() + filepath.Ext(header.Filename)
	err = bucket.PutObject(key, file)
	if err != nil {
		logx.Errorf("PutObject error: %v", err)
		return nil, ecode.PutBucketObjectErr
	}
	return &types.UploadCoverResponse{CoverUrl: fileUrl(key)}, nil
}

func fileUrl(key string) string {
	return "https://nddwl.oss-cn-beijing.aliyuncs.com/" + key
}
