package logic

import (
	"Cloud-Disk/core/helper"
	"context"

	"Cloud-Disk/core/internal/svc"
	"Cloud-Disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAuthLogic {
	return &RefreshAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshAuthLogic) RefreshAuth(req *types.RefreshAuthRequest, authorization string) (resp *types.RefreshAuthResponse, err error) {
	//// 解析 Authorization 获取 UserClaim
	userClaim, err := helper.AnalyseToken(authorization)
	if err != nil {
		return
	}
	//生成新的token和RefreshToken
	token, err := helper.GenerateToken(userClaim.ID, userClaim.Identity, userClaim.Name, 3600)
	if err != nil {
		return
	}
	refreshToken, err := helper.GenerateToken(userClaim.ID, userClaim.Identity, userClaim.Name, 3600)
	if err != nil {
		return
	}
	resp = new(types.RefreshAuthResponse)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
