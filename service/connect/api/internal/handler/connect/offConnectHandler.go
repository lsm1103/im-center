package connect

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/common/result"
	"im-center/service/connect/api/internal/logic/connect"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"
)

func OffConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConnectUid
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.Error(w, err)
			result.ParamErrorResult(r, w, err)
			return
		}

		l := connect.NewOffConnectLogic(r.Context(), svcCtx)
		resp, err := l.OffConnect(req)
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
