package logic

import (
	"Cloud-Disk/core/models"
	"context"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareDetailLogic {
	return &ShareDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareDetailLogic) ShareDetail(req *types.ShareDetailRequest) (resp *types.ShareDetailResponse, err error) {
	// 对分享记录的点击次数进行 + 1
	_, err = models.Engine.Exec("UPDATE share SET click_num = click_num + 1 WHERE identity = ?", req.Identity)
	if err != nil {
		return
	}
	// 获取文件的详细信息
	resp = new(types.ShareDetailResponse)
	_, err = models.Engine.Table("share").
		Select("share.repository_identity,user_repository.name,repository_pool.ext,repository_pool.size,repository_pool.path").
		Join("LEFT", "user_repository", "share.user_repository_identity = user_repository.identity").
		Join("LEFT", "repository_pool", "share.repository_identity=repository_pool.identity").
		Where("share.identity = ? ", req.Identity).Get(resp)
	return
}
