// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"lark/apps/ws/internal/config"
	"lark/apps/ws/internal/server"

	"github.com/gorilla/websocket"
)

type ServiceContext struct {
	Config config.Config
	MsgServer *server.Server
	WsUpgrader *websocket.Upgrader
}

func NewServiceContext(
	c config.Config,
	msgserver *server.Server,
	wsupgrader *websocket.Upgrader) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MsgServer: msgserver,
		WsUpgrader: wsupgrader,
	}
}
