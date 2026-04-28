package logic

import (
	"errors"
	"strings"
	"time"

	"lark/apps/user/api/internal/types"
	"lark/apps/user/rpc/userrpc"
	"lark/pkg/utils"
)

func toAPIUserProfile(profile *userrpc.UserProfile) types.UserProfile {
	if profile == nil {
		return types.UserProfile{}
	}

	createdAt := ""
	if profile.CreatedAt > 0 {
		createdAt = time.Unix(profile.CreatedAt, 0).Format(time.RFC3339)
	}

	updatedAt := ""
	if profile.UpdatedAt > 0 {
		updatedAt = time.Unix(profile.UpdatedAt, 0).Format(time.RFC3339)
	}

	return types.UserProfile{
		Uid:       profile.Uid,
		Name:      profile.Name,
		Email:     profile.Email,
		Phone:     profile.Phone,
		Status:    int8(profile.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func parseClaimsFromAuthorizationHeader(authHeader, secret string) (*utils.Claims, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid Authorization header, expected: Bearer <token>")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return nil, errors.New("empty bearer token")
	}

	return utils.ParseToken(secret, token)
}
