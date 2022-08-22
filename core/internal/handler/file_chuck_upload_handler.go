package handler

import (
	"Cloud-Disk/core/helper"
	"errors"
	"net/http"

	"Cloud-Disk/core/internal/logic"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileChuckUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileChuckUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		//检查参数是否填写
		if r.PostForm.Get("key") == "" {
			httpx.Error(w, errors.New("key为空"))
			return
		}
		if r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("upload_id为空"))
			return
		}
		if r.PostForm.Get("part_number") == "" {
			httpx.Error(w, errors.New("part_number为空"))
			return
		}
		etag, err := helper.CosFileChuck(r)
		if err != nil {
			httpx.Error(w, err)
		}

		l := logic.NewFileChuckUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileChuckUpload(&req)
		resp = new(types.FileChuckUploadResponse)
		resp.ETag = etag
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
