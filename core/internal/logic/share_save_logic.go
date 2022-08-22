package logic

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"Cloud-Disk/core/models"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareSaveLogic {
	return &ShareSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareSaveLogic) ShareSave(req *types.ShareSaveRequest, userIdentity string) (resp *types.ShareSaveResponse, err error) {
	//根据 repository_identity 获取资源详情
	rp := new(models.RepositoryPool)
	has, err := models.Engine.Where("identity= ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("该资源不存在！")
	}
	//保存资源到个人存储池 user_repository
	ur := &models.UserRepository{
		Identity:           helper.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Name:               rp.Name,
		Ext:                rp.Ext,
	}
	_, err = models.Engine.Insert(ur)
	resp = new(types.ShareSaveResponse)
	resp.Identity = ur.Identity
	return
}
