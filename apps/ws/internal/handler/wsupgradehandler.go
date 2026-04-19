// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lark/apps/ws/internal/logic"
	"lark/apps/ws/internal/svc"
	"lark/apps/ws/internal/types"
)

func WsUpgradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewWsUpgradeLogic(r.Context(), svcCtx)
		resp, err := l.WsUpgrade(&req,w,r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		if resp == nil {
			return // 升级成功后不返回任何响应
		}
	}
}
