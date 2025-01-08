package handler

import (
	"applet/internal/ecode"
	"net/http"

	"applet/internal/logic"
	"applet/internal/svc"
	"applet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VerificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerificationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewVerificationLogic(r.Context(), svcCtx)
		resp, err := l.Verification(&req)
		ecode.JsonCtx(r.Context(), w, resp, err)
	}
}
