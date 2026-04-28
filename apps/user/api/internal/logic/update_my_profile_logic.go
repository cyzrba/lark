package logic

import (
	"context"
	"errors"

	"lark/apps/user/api/internal/svc"
	"lark/apps/user/api/internal/types"
	"lark/apps/user/rpc/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyProfileLogic {
	return &UpdateMyProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyProfileLogic) UpdateMyProfile(req *types.UpdateMyProfileRequest) (resp *types.UpdateMyProfileResponse, err error) {
	claims, err := parseClaimsFromAuthorizationHeader(req.Authorization, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.UserRpc.UpdateUser(l.ctx, &userrpc.UpdateUserReq{
		Uid:          claims.Uid,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       int32(req.Status),
		UpdateStatus: req.UpdateStatus,
		OldPassword:  req.OldPassword,
		NewPassword:  req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil || rpcResp.User == nil {
		return nil, errors.New("empty update response from rpc")
	}

	return &types.UpdateMyProfileResponse{User: toAPIUserProfile(rpcResp.User)}, nil
}
