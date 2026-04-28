// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1
//goctl api go -api <你的文件名>.api -dir . -style go_zero
package main

import (
	"flag"
	"fmt"
	"net/http"

	"lark/apps/ws/internal/config"
	"lark/apps/ws/internal/handler"
	"lark/apps/ws/internal/svc"
	msgserver "lark/apps/ws/internal/msgserver"
	"lark/apps/chat/rpc/chatrpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/zrpc"
	
)

var configFile = flag.String("f", "etc/ws-api.yaml", "the config file")

func main() {
	flag.Parse()
	//反序列化配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	
	ChatRpc := chatrpc.NewChatRpc(zrpc.MustNewClient(c.ChatRpc))
	//创建rest server
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	//创建消息服务器
	msgServer := msgserver.NewServer(ChatRpc)
	fmt.Printf("msgServer started at %s:%d.\n", c.Host, c.Port)
	//todo: 如何优雅关闭？
	go msgServer.Start()


	//创建service context（服务上下文）
	ctx := svc.NewServiceContext(c,msgServer,&websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},ChatRpc)
	//注册路由
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
	
}
