package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Driver       string `json:"driver"`
	Source       string `json:"source"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}
