package logic

import (
	"errors"
	"strings"

	"lark/apps/user/model"
	"lark/apps/user/rpc/types/user"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

const mysqlDuplicateEntryErrCode uint16 = 1062

func mapUserProfile(u *model.User) *user.UserProfile {
	if u == nil {
		return nil
	}

	return &user.UserProfile{
		Uid:       u.Uuid,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Status:    int32(u.Status),
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}

func normalizeStr(s string) string {
	return strings.TrimSpace(s)
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == mysqlDuplicateEntryErrCode
	}
	return false
}
