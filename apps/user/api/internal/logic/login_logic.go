package logic

import (
	"context"
	"errors"

	"lark/apps/user/api/internal/svc"
	"lark/apps/user/api/internal/types"
	"lark/apps/user/rpc/userrpc"
	"lark/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &userrpc.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil || rpcResp.User == nil {
		return nil, errors.New("empty login response from rpc")
	}

	token, expiresAt, err := utils.GenerateToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
		rpcResp.User.Uid,
		rpcResp.User.Name,
	)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
