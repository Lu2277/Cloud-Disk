package handler

import (
	"net/http"

	"Cloud-Disk/core/internal/logic"
	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShareRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShareRecordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewShareRecordLogic(r.Context(), svcCtx)
		resp, err := l.ShareRecord(&req, r.Header.Get("userIdentity"))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
