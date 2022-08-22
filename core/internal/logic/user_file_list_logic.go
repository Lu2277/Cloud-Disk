package logic

import (
	"Cloud-Disk/core/models"
	"context"
	"time"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, userIdentity string) (resp *types.UserFileListResponse, err error) {
	//设置分页参数
	size := req.Size
	if size <= 0 {
		size = 10 //默认10条数据
	}
	page := req.Page
	if page <= 0 {
		page = 1 //默认1页
	}
	offset := (page - 1) * size
	//在数据库中查询用户文件列表(软删除的数据也显示出来)
	userFileList := make([]*types.UserFile, 0)
	err = models.Engine.Table("user_repository").Where("user_identity=? AND parent_id=?", userIdentity, req.Id).
		Select("user_repository.id,user_repository.identity,user_repository.repository_identity,"+
			"user_repository.name,user_repository.ext,repository_pool.size,repository_pool.path").
		Join("LEFT", "repository_pool", "user_repository.repository_identity=repository_pool.identity").
		Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format("2006-01-02 15:04:05")).
		Limit(size, offset).Find(&userFileList)
	if err != nil {
		return
	}
	//查询文件总数
	count, err := models.Engine.Where("user_identity=? AND parent_id=?", userIdentity, req.Id).Count(new(models.UserRepository))
	if err != nil {
		return
	}
	resp = new(types.UserFileListResponse)
	resp.List = userFileList
	resp.Count = count
	return
}
