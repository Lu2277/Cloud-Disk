package logic

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"Cloud-Disk/core/models"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest) (resp *types.FileUploadResponse, err error) {
	rp := &models.RepositoryPool{
		Identity: helper.GetUUID(),
		Name:     req.Name,
		Hash:     req.Hash,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	_, err = models.Engine.Insert(rp)
	if err != nil {
		return nil, err
	}
	resp = new(types.FileUploadResponse)
	resp.Identity = rp.Identity
	resp.Name = rp.Name
	resp.Ext = rp.Ext
	return
}
