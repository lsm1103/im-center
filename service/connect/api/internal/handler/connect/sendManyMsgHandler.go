package connect

import (
	"encoding/json"
	"net/http"

	"im-center/common/result"
	"im-center/service/connect/api/internal/logic/connect"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"
)

func SendManyMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendManyMsgReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	// httpx.Error(w, err)
		//	result.ParamErrorResult(r, w, err)
		//	return
		//}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := connect.NewSendManyMsgLogic(r.Context(), svcCtx)
		resp, err := l.SendManyMsg(req)
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
