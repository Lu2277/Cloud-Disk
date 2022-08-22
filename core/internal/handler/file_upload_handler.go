package handler

import (
	"Cloud-Disk/core/helper"
	"Cloud-Disk/core/models"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"

	"Cloud-Disk/core/internal/logic"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		//获取文件
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			return
		}
		bytes := make([]byte, fileHeader.Size)
		_, err = file.Read(bytes)
		if err != nil {
			return
		}
		//获取文件的hash值
		hash := fmt.Sprintf("%x", md5.Sum(bytes))
		rp := new(models.RepositoryPool)
		//判断文件是否已存在
		has, err := models.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			return
		}
		if has {
			//文件存在，返回文件的identity、ext、name信息
			httpx.OkJson(w, &types.FileUploadResponse{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		//文件不存在，上传文件到腾讯云
		uploadPath, err := helper.CosUpload(r)
		if err != nil {
			return
		}
		//往logic层传递req信息
		req.Name = fileHeader.Filename
		req.Hash = hash
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Path = uploadPath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
