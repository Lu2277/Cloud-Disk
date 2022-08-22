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

type ShareRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareRecordLogic {
	return &ShareRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareRecordLogic) ShareRecord(req *types.ShareRecordRequest, userIdentity string) (resp *types.ShareRecordResponse, err error) {
	// 判断用户存储池中是否存在该分享文件
	ur := new(models.UserRepository)
	has, err := models.Engine.Where("identity = ? AND user_identity = ?",
		req.UserRepositoryIdentity, userIdentity).Get(ur)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("该分享文件不存在！")
	}
	shareData := &models.Share{
		Identity:               helper.GetUUID(),
		UserIdentity:           userIdentity,
		RepositoryIdentity:     ur.RepositoryIdentity,
		UserRepositoryIdentity: req.UserRepositoryIdentity,
		ExpiredTime:            req.ExpiredTime,
	}
	_, err = models.Engine.Insert(shareData)
	if err != nil {
		return
	}
	resp = new(types.ShareRecordResponse)
	resp.Identity = shareData.Identity
	return
}
