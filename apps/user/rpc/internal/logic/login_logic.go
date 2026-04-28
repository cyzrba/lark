package logic

import (
	"context"
	"errors"

	"lark/apps/user/model"
	"lark/apps/user/rpc/internal/svc"
	"lark/apps/user/rpc/types/user"
	"lark/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	email := normalizeStr(in.Email)
	if email == "" || in.Password == "" {
		return nil, errors.New("email and password are required")
	}

	var entity model.User
	err := l.svcCtx.DB.WithContext(l.ctx).Where("email = ?", email).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.VerifyPassword(entity.Password, in.Password) {
		return nil, errors.New("invalid email or password")
	}

	return &user.LoginResp{User: mapUserProfile(&entity)}, nil
}
