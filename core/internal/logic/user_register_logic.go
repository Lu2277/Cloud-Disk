package logic

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/models"
	"context"
	"errors"
	"log"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	//判断验证码是否正确
	code, err := models.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("获取不到该邮箱的验证码,请检查邮箱是否正确")
	}
	if code != req.Code {
		err = errors.New("验证码不正确")
		return
	}
	//判断用户名是否存在
	count, err := models.Engine.Where("name = ?", req.Name).Count(new(models.User))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		err = errors.New("该用户名已存在")
	}
	//都没有问题，将用户信息插入数据库
	user := &models.User{
		Identity: helper.GetUUID(),
		Name:     req.Name,
		Password: helper.GetMd5(req.Password),
		Email:    req.Email,
	}
	n, err := models.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	log.Println("Success！Insert User Row:", n)
	resp = new(types.RegisterResponse)
	resp.Message = "注册成功！"
	return
}
