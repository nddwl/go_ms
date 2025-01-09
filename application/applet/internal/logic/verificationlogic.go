package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"
	"zhihu/application/applet/internal/svc"
	"zhihu/application/applet/internal/types"
	"zhihu/application/user/user"
	"zhihu/pkg/ecode"
	"zhihu/pkg/utils"
	"zhihu/pkg/validator"
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	if err = validator.Struct(req); err != nil {
		return nil, ecode.RequestErr
	}
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error %v", req.Mobile, err)
		return nil, err
	}
	if count > types.VerificationLimitPerDay {
		return nil, ecode.VerificationMaxLimit
	}
	code := utils.GenerateVerificationCode()
	_, err = l.svcCtx.UserRpc.SendSms(l.ctx, &user.SendSmsRequest{Mobile: req.Mobile})
	if err != nil {
		logx.Errorf("UserRpc->sendSms mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = l.saveVerificationCode(req.Mobile, code)
	if err != nil {
		logx.Errorf("saveVerificationCode mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerification mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	count, err := l.svcCtx.BizRedis.Get(types.Prefix + ":count:" + mobile)
	if err != nil {
		return 0, err
	}
	if count == "" {
		return 0, nil
	}
	return strconv.Atoi(count)
}

func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := types.Prefix + ":count:" + mobile
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return nil
	}
	return l.svcCtx.BizRedis.Expireat(key, utils.EndOfDay(time.Now()).Unix())
}

func (l *VerificationLogic) saveVerificationCode(mobile string, code string) error {
	return l.svcCtx.BizRedis.Setex(types.Prefix+":code:"+mobile, code, types.ExpireTime)
}

func verifyVerificationCode(redis *redis.Redis, mobile string, code string) (bool, error) {
	c, err := redis.Get(types.Prefix + ":code:" + mobile)
	switch {
	case err != nil:
		return false, err
	case c != code:
		return false, nil
	default:
		return true, nil
	}
}

func delVerificationCode(redis *redis.Redis, mobile string) error {
	_, err := redis.Del(types.Prefix + ":code:" + mobile)
	return err
}
