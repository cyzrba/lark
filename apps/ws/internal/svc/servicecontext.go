// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"lark/apps/ws/internal/config"
	"lark/apps/ws/internal/msgserver"
	"lark/apps/chat/rpc/chatrpc"

	"github.com/gorilla/websocket"
)

type ServiceContext struct {
	Config config.Config
	MsgServer *server.Server
	WsUpgrader *websocket.Upgrader
	ChatRpc chatrpc.ChatRpc
}

func NewServiceContext(
	c config.Config,
	msgserver *server.Server,
	wsupgrader *websocket.Upgrader,
	chatRpc chatrpc.ChatRpc) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MsgServer: msgserver,
		WsUpgrader: wsupgrader,
		ChatRpc: chatRpc,
	}
}
