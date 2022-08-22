package handler

import (
	"net/http"

	"Cloud-Disk/core/internal/logic"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshAuthRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRefreshAuthLogic(r.Context(), svcCtx)
		resp, err := l.RefreshAuth(&req, r.Header.Get("Authorization")) //往logic层传递需要的Authorization
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
