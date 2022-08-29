package friend

import (
	"net/http"

	//"pj-auxiliary-diagnosis-C/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/service/business/chatService/api/internal/logic/friend"
	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"
)

func GroupGetsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetsReq
		if err := httpx.Parse(r, &req); err != nil {
			// httpx.Error(w, err)
			result.ParamErrorResult(r, w, err)
			return
		}

		l := friend.NewGroupGetsLogic(r.Context(), svcCtx)
		resp, err := l.GroupGets(req)
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
