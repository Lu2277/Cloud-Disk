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

type FileUploadPrepareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPrepareLogic {
	return &FileUploadPrepareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest) (resp *types.FileUploadPrepareResponse, err error) {
	rp := new(models.RepositoryPool)
	//在数据库中查找有无该文件（查md5值）
	has, err := models.Engine.Where("hash = ?", req.Hash).Get(rp)
	if err != nil {
		return
	}
	resp = new(types.FileUploadPrepareResponse)
	if has {
		//存储池中已存在该文件，则直接秒传
		resp.Identity = rp.Identity
	} else {
		//	获取uploadId和key，用于进行分片上传处理
		key, upLoadId, err := helper.CosFileChuckPrepare(req.Ext)
		if err != nil {
			return nil, errors.New("分片预处理失败！")
		}
		resp.UploadId = upLoadId
		resp.Key = key
	}
	return
}
