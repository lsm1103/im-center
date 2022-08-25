package result

import (
	"net/http"

	"im-center/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

//http返回
func HttpResult(r *http.Request, w http.ResponseWriter, err error, resp interface{}) {
	if err == nil {
		//成功返回
		httpx.WriteJson(w, http.StatusOK, Success(resp))
	} else {
		//错误返回
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"
		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := gstatus.Code()
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}
		logx.WithContext(r.Context()).Errorf("【%s-API-ERR】】:%s, resp:%+v", r.RequestURI, err.Error(), resp )
		if errcode == xerr.OK{
			httpx.OkJson(w, &ResponseSuccessBean{
				Code: errcode,
				Msg:  errmsg,
				Data: resp,
			})
			return
		} else if errcode == xerr.ALREADY_EXISTS{
			httpx.OkJson(w, &ResponseSuccessBean{
				Code: errcode,
				Msg:  errmsg,
				Data: resp,
			})
			return
		}
		httpx.WriteJson(w, http.StatusBadRequest, Error(errcode, errmsg))
	}
}

//授权的http方法
func AuthHttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"
		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := gstatus.Code()
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}
		logx.WithContext(r.Context()).Errorf("【%s-GATEWAY-AuthHttpResult-ERR】:%s, resp:%+v", r.RequestURI, err.Error(), resp)
		httpx.WriteJson(w, http.StatusUnauthorized, Error(errcode, errmsg))
	}
}

//http 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	code := xerr.REUQES_PARAM_ERROR
	logx.WithContext(r.Context()).Errorf("【%s-GATEWAY-ParamErrorResult-ERR】:%s, ERR:%+v", r.RequestURI, xerr.MapErrMsg(code), err.Error() )
	httpx.WriteJson(w, http.StatusBadRequest, Error(code, xerr.MapErrMsg(code) ) )
}

//http 请求验证错误返回
func RequestVerifyErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	if err != nil {
		code := xerr.REQUEST_VERIFY_ERROR
		logx.WithContext(r.Context()).Errorf("【%s-GATEWAY-RequestVerify-ERR】:%s, ERR:%+v", r.RequestURI, xerr.MapErrMsg(code), err.Error() )
		httpx.WriteJson(w, http.StatusBadRequest, Error(code, xerr.MapErrMsg(code) ) )
	}
}
