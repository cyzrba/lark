package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Uid  int64  `json:"uid"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, expireSeconds int64, uid int64, name string) (string, int64, error) {
	if secret == "" {
		return "", 0, errors.New("jwt secret is empty")
	}
	if expireSeconds <= 0 {
		return "", 0, errors.New("jwt expire seconds must be greater than zero")
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(expireSeconds) * time.Second).Unix()
	claims := Claims{
		Uid:  uid,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return signed, expiresAt, nil
}

func ParseToken(secret, tokenString string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
