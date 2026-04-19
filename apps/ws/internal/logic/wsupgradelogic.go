// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"net/http"

	"lark/apps/ws/internal/svc"
	"lark/apps/ws/internal/types"
	"lark/apps/ws/internal/client"

	"github.com/zeromicro/go-zero/core/logx"
)

type WsUpgradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWsUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsUpgradeLogic {
	return &WsUpgradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WsUpgradeLogic) WsUpgrade(req *types.Request, w http.ResponseWriter, r *http.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	conn, err := l.svcCtx.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := client.NewClient(
		req.Name,
		conn, 
		l.svcCtx.MsgServer,
	)
	l.svcCtx.MsgServer.AddClient(c)
	go c.Start()
	return nil, nil
}
