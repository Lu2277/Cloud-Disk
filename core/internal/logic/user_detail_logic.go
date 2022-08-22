package logic

import (
	"Cloud-Disk/core/models"
	"context"
	"errors"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailRequest) (resp *types.UserDetailResponse, err error) {
	resp = &types.UserDetailResponse{}
	user := new(models.User)
	has, err := models.Engine.Where("identity=?", req.Identity).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("没有该用户")
	}
	resp.Name = user.Name
	resp.Email = user.Email
	return
}
