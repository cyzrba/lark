package logic

import (
	"context"
	"errors"

	"lark/apps/user/model"
	"lark/apps/user/rpc/internal/svc"
	"lark/apps/user/rpc/types/user"
	"lark/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	name := normalizeStr(in.Name)
	email := normalizeStr(in.Email)
	phone := normalizeStr(in.Phone)
	if name == "" || email == "" || in.Password == "" {
		return nil, errors.New("name, email and password are required")
	}

	hashedPassword, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	uid, err := utils.NextSnowflakeID()
	if err != nil {
		return nil, err
	}

	entity := &model.User{
		Uuid:     uid,
		Name:     name,
		Password: hashedPassword,
		Email:    email,
		Phone:    phone,
		Status:   1,
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Create(entity).Error; err != nil {
		if isDuplicateEntryError(err) {
			return nil, errors.New("email or phone already exists")
		}
		return nil, err
	}

	return &user.RegisterResp{User: mapUserProfile(entity)}, nil
}
