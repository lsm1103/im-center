package result

import "google.golang.org/grpc/codes"

type ResponseSuccessBean struct {
	Code codes.Code      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type NullJson struct{}

func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{200, "OK", data}
}

type ResponseErrorBean struct {
	Code codes.Code `json:"code"`
	Msg  string `json:"msg"`
}

func Error(errCode codes.Code, errMsg string) *ResponseErrorBean {
	return &ResponseErrorBean{errCode, errMsg}
}
