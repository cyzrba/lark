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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserReq) (*user.UpdateUserResp, error) {
	if in.Uid <= 0 {
		return nil, errors.New("invalid uid")
	}

	db := l.svcCtx.DB.WithContext(l.ctx)

	var entity model.User
	if err := db.Where("uid = ?", in.Uid).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})
	if name := normalizeStr(in.Name); name != "" {
		updates["name"] = name
	}
	if email := normalizeStr(in.Email); email != "" {
		updates["email"] = email
	}
	if phone := normalizeStr(in.Phone); phone != "" {
		updates["phone"] = phone
	}
	if in.UpdateStatus {
		updates["status"] = int8(in.Status)
	}

	if in.NewPassword != "" {
		if in.OldPassword == "" {
			return nil, errors.New("old password is required when changing password")
		}
		if !utils.VerifyPassword(entity.Password, in.OldPassword) {
			return nil, errors.New("old password is incorrect")
		}
		hashedPassword, err := utils.HashPassword(in.NewPassword)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashedPassword
	}

	if len(updates) > 0 {
		if err := db.Model(&entity).Updates(updates).Error; err != nil {
			if isDuplicateEntryError(err) {
				return nil, errors.New("email or phone already exists")
			}
			return nil, err
		}
	}

	if err := db.Where("uid = ?", in.Uid).First(&entity).Error; err != nil {
		return nil, err
	}

	return &user.UpdateUserResp{User: mapUserProfile(&entity)}, nil
}
