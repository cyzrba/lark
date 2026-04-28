package logic

import (
	"context"

	"lark/apps/chat/rpc/internal/svc"
	"lark/pkg/proto/pb/chat"//todo

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


func (l *SendMessageLogic) SendMessage(in *pb_chat.SendMessageReq) (*pb_chat.SendMessageResp, error) {
	return &pb_chat.SendMessageResp{}, nil
} 