package logic

import (
	"Cloud-Disk/core/models"
	"context"
	"errors"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest, userIdentity string) (resp *types.UserFileNameUpdateResponse, err error) {
	data := &models.UserRepository{Name: req.Name}
	// 判断当前名称在该层级下是否存在
	count, err := models.Engine.Where("name=? AND parent_id = (SELECT parent_id FROM user_repository  WHERE user_repository.identity= ?)",
		req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		return
	}
	if count > 0 {
		return nil, errors.New("该文件名已存在")
	}
	//文件名称修改
	_, err = models.Engine.Where("identity=? AND user_identity=?", req.Identity, userIdentity).Update(data)
	if err != nil {
		return nil, err
	}
	return
}
