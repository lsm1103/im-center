package msg

import (
	"net/http"

	//"pj-auxiliary-diagnosis-C/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/service/business/chatService/api/internal/logic/msg"
	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"
)

func BatchDelMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDelMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.Error(w, err)
			result.ParamErrorResult(r, w, err)
			return
		}

		l := msg.NewBatchDelMsgLogic(r.Context(), svcCtx)
		err := l.BatchDelMsg(req)
		/*
			if err != nil {
				httpx.Error(w, err)
			} else {
				httpx.Ok(w)
			}
		*/
		result.HttpResult(r, w, err, nil)
	}
}
