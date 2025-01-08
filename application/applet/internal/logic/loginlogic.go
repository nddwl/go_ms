package logic

import (
	"context"
	"zhihu/application/applet/service"
	"zhihu/pkg/ecode"
	"zhihu/pkg/utils"

	"zhihu/application/applet/internal/svc"
	"zhihu/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	if err = l.svcCtx.Validator.Struct(req); err != nil {
		err = ecode.BadRequest
		return
	}
	ok, err := verifyVerificationCode(l.svcCtx.Redis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("verifyVerificationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if !ok {
		return nil, ecode.VerificationCodeFailed
	}
	mobileResp, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &service.FindByMobileRequest{Mobile: req.Mobile})
	if err != nil {
		logx.Errorf("userRpc->FindByMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if mobileResp == nil && mobileResp.UserId == 0 {
		return nil, ecode.UserNotExisted
	}
	auth := l.svcCtx.Config.Auth
	token, err := utils.GenerateToken(auth.AccessSecret, auth.AccessExpire, map[string]interface{}{"userId": mobileResp.UserId})
	if err != nil {
		logx.Errorf("GenerateToken userId: %d error: %v", mobileResp.UserId, err)
		return
	}
	return &types.LoginResponse{
		UserId: mobileResp.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}
