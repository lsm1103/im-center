package service

import (
	"net/http"

	//"pj-auxiliary-diagnosis-C/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/service/business/chatService/api/internal/logic/service"
	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"
)

func OffConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConnectUid
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.Error(w, err)
			result.ParamErrorResult(r, w, err)
			return
		}

		l := service.NewOffConnectLogic(r.Context(), svcCtx)
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
