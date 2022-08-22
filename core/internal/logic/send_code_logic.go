package logic

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"Cloud-Disk/core/models"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type SendCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCodeLogic) SendCode(req *types.SendCodeRequest) (resp *types.SendCodeResponse, err error) {
	user := new(models.User)
	i, err := models.Engine.Where("email = ?", req.Email).Count(user)
	if err != nil {
		return
	}
	if i > 0 {
		err = errors.New("该邮箱已注册")
		return
	}
	//设置验证码有效期
	models.RDB.Set(l.ctx, req.Email, helper.CreateRandCode(), 300*time.Second)
	// 发送随机验证码
	err = helper.SendCode(req.Email, helper.CreateRandCode())
	return
}
