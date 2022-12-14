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

type UserFolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest, userIdentity string) (resp *types.UserFolderCreateResponse, err error) {
	// 判断当前名称在该层级下是否存在
	count, err := models.Engine.Where("name = ? AND parent_id = ?", req.Name, req.ParentId).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该名称已存在")
	}
	data := &models.UserRepository{
		Identity:     helper.GetUUID(),
		UserIdentity: userIdentity,
		ParentId:     req.ParentId,
		Name:         req.Name,
	}
	_, err = models.Engine.Insert(data)
	if err != nil {
		return
	}
	resp = new(types.UserFolderCreateResponse)
	resp.Identity = data.Identity
	return
}
