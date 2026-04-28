package logic

import (
	"context"
	"errors"

	"lark/apps/user/api/internal/svc"
	"lark/apps/user/api/internal/types"
	"lark/apps/user/rpc/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyProfileLogic {
	return &GetMyProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyProfileLogic) GetMyProfile(req *types.GetMyProfileRequest) (resp *types.GetMyProfileResponse, err error) {
	claims, err := parseClaimsFromAuthorizationHeader(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.UserRpc.GetUserDetail(l.ctx, &userrpc.GetUserDetailReq{Uid: claims.Uid})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil || rpcResp.User == nil {
		return nil, errors.New("empty profile response from rpc")
	}

	return &types.GetMyProfileResponse{User: toAPIUserProfile(rpcResp.User)}, nil
}
