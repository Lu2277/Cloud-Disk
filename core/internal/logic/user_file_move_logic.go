package logic

import (
	"Cloud-Disk/core/models"
	"context"
	"errors"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest, userIdentity string) (resp *types.UserFileMoveResponse, err error) {
	// 判断父级文件夹是否存在
	parentData := new(models.UserRepository)
	has, err := models.Engine.Where("identity = ? AND user_identity = ?", req.ParentIdentity, userIdentity).Get(parentData)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("文件夹不存在")
	}
	//更新传入的parent_id
	_, err = models.Engine.Where("identity = ?", req.Identity).Update(models.UserRepository{
		ParentId: int64(parentData.Id),
	})
	return
}
