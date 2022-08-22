package logic

import (
	"Cloud-Disk/core/helper"
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileChuckUploadCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileChuckUploadCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileChuckUploadCompleteLogic {
	return &FileChuckUploadCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileChuckUploadCompleteLogic) FileChuckUploadComplete(req *types.FileChuckUploadCompleteRequest) (resp *types.FileChuckUploadCompleteResponse, err error) {
	co := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		co = append(co, cos.Object{
			PartNumber: v.PartNumber,
			ETag:       v.Etag,
		})
	}
	err = helper.CosCompleteChuck(req.Key, req.UploadId, co)
	return
}
