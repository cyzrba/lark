package logic

import (
	"context"
	"errors"

	"lark/apps/user/api/internal/svc"
	"lark/apps/user/api/internal/types"
	"lark/apps/user/rpc/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Register(l.ctx, &userrpc.RegisterReq{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil || rpcResp.User == nil {
		return nil, errors.New("empty register response from rpc")
	}

	return &types.RegisterResponse{
		User: toAPIUserProfile(rpcResp.User),
	}, nil
}
