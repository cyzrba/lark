package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpc zrpc.RpcClientConf `json:"userRpc"`
	JwtAuth JwtAuthConfig      `json:"jwtAuth"`
}

type JwtAuthConfig struct {
	AccessSecret string `json:"access_secret"`
	AccessExpire int64  `json:"access_expire"`
}
