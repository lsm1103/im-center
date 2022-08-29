package msg

import (
	"net/http"

	//"pj-auxiliary-diagnosis-C/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/service/business/chatService/api/internal/logic/msg"
	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"
)

func SyncMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SyncMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.Error(w, err)
			result.ParamErrorResult(r, w, err)
			return
		}

		l := msg.NewSyncMsgLogic(r.Context(), svcCtx)
		resp, err := l.SyncMsg(req)
		/*
			if err != nil {
				httpx.Error(w, err)
			} else {
				httpx.OkJson(w, resp)
			}
		*/
		result.HttpResult(r, w, err, resp)
	}
}
