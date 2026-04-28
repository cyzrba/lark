package logic

import (
	"context"
	"errors"

	"lark/apps/user/model"
	"lark/apps/user/rpc/internal/svc"
	"lark/apps/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserDetailLogic) GetUserDetail(in *user.GetUserDetailReq) (*user.GetUserDetailResp, error) {
	if in.Uid <= 0 {
		return nil, errors.New("invalid uid")
	}

	var entity model.User
	if err := l.svcCtx.DB.WithContext(l.ctx).Where("uid = ?", in.Uid).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user.GetUserDetailResp{User: mapUserProfile(&entity)}, nil
}
