package handler

import (
	"net/http"

	"Cloud-Disk/core/internal/logic"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserFolderCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFolderCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUserFolderCreateLogic(r.Context(), svcCtx)
		resp, err := l.UserFolderCreate(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
