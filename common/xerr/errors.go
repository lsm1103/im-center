package xerr

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
常用通用固定错误
*/

type CodeError struct {
	Code codes.Code    `json:"code"`
	Msg  string    `json:"msg"`
}

//返回给前端的错误码
func (e *CodeError) GetErrCode() codes.Code {
	return e.Code
}

//返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.Msg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("Code:%d，Msg:%s;", e.Code, e.Msg)
}

func NewErrCodeMsg(Code codes.Code, Msg string) *CodeError {
	return &CodeError{Code: Code, Msg: Msg}
}
func NewErrCode(Code codes.Code) *CodeError {
	return &CodeError{Code: Code, Msg: MapErrMsg(Code)}
}

func NewErrMsg(Msg string) *CodeError {
	return &CodeError{Code: SERVER_COMMON_ERROR, Msg: Msg}
}

func NewErrMsgf(Msgf string, arg ...interface{}) *CodeError {
	return &CodeError{Code: SERVER_COMMON_ERROR, Msg: fmt.Sprintf(Msgf, arg...)}
}

type RpcCodeError struct {
	logx.Logger
}

func NewRpcCodeError(log logx.Logger) *RpcCodeError {
	return &RpcCodeError{log }
}

func (r RpcCodeError)NewRpcErrCode(Code codes.Code, RealMsg string) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", RealMsg )
	return status.Error(Code, MapErrMsg(Code) )
}

func (r RpcCodeError)NewRpcErrCodef(Code codes.Code, RealMsgf string, rArg ...interface{}) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", fmt.Sprintf(RealMsgf, rArg...) )
	return status.Error(Code, MapErrMsg(Code) )
}

func (r RpcCodeError)NewRpcErrMsg(Msg string, RealMsg string) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", RealMsg )
	return status.Error(SERVER_COMMON_ERROR, Msg)
}

func (r RpcCodeError)NewRpcErrMsgf(Msg string, RealMsgf string, rArg ...interface{}) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", fmt.Sprintf(RealMsgf, rArg...) )
	return status.Error(SERVER_COMMON_ERROR, Msg)
}

func (r RpcCodeError)NewRpcErrCodeMsg(Code codes.Code, Msg string, RealMsg string) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", RealMsg )
	return status.Error(Code, Msg)
}

func (r RpcCodeError)NewRpcErrCodeMsgf(Code codes.Code, Msg string, RealMsgf string, rArg ...interface{}) error {
	r.Logger.Errorf("【RPC-ERR】】:%s", fmt.Sprintf(RealMsgf, rArg...) )
	return status.Error(Code, Msg)
}
