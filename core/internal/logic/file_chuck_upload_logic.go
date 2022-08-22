package logic

import (
	"context"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileChuckUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileChuckUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileChuckUploadLogic {
	return &FileChuckUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileChuckUploadLogic) FileChuckUpload(req *types.FileChuckUploadRequest) (resp *types.FileChuckUploadResponse, err error) {
	return
}
