package logic

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/models"
	"context"
	"errors"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	//1.从数据查询当前用户
	user := new(models.User)
	has, err := models.Engine.Where("name=? AND password=?", req.Name, helper.GetMd5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名或者密码错误")
	}
	//2.如果查询到该用户，返回用户登录的token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, 3600)
	if err != nil {
		return nil, errors.New("生成token错误")
	}
	//生成用于刷新token的RefreshToken
	refreshToken, err := helper.GenerateToken(user.Id, user.Identity, user.Name, 7200)
	if err != nil {
		return nil, errors.New("生成RefreshToken错误")
	}
	resp = new(types.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
