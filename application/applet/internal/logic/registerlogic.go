package logic

import (
	"context"
	"zhihu/application/applet/internal/code"
	"zhihu/application/applet/internal/svc"
	"zhihu/application/applet/internal/types"
	"zhihu/application/user/user"
	"zhihu/pkg/ecode"
	"zhihu/pkg/utils"
	"zhihu/pkg/validator"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	if err = validator.Struct(req); err != nil {
		return nil, ecode.RequestErr
	}
	ok, err := verifyVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("verifyVerificationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if !ok {
		return nil, code.VerificationCodeFailed
	}
	mobileResp, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &user.FindByMobileRequest{})
	if err != nil {
		logx.Errorf("UserRpc->FindByMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if mobileResp.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	registerResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   req.Mobile,
		Avatar:   "",
		Password: utils.GenerateFromPassword(req.Password),
	})
	if err != nil {
		logx.Errorf("UserRc->Register req: %s error: %v", req, err)
		return nil, err
	}
	auth := l.svcCtx.Config.Auth
	token, err := utils.GenerateToken(auth.AccessSecret, auth.AccessExpire, map[string]interface{}{"userId": registerResp.UserId})
	if err != nil {
		logx.Errorf("GenerateToken userId: %d error: %v", registerResp.UserId, err)
		return
	}
	err = delVerificationCode(l.svcCtx.BizRedis, req.Mobile)
	if err != nil {
		logx.Errorf("delVerificationCode error: %v", err)
	}
	return &types.RegisterResponse{
		UserId: registerResp.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}
