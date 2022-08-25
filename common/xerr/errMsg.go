package xerr

import "google.golang.org/grpc/codes"

var message map[codes.Code]string

func init() {
	message = make(map[codes.Code]string)
	message[OK] = "SUCCESS"
	message[SERVER_COMMON_ERROR] = "服务器开小差啦,稍后再来试一试"
	message[REUQES_PARAM_ERROR] = "参数错误"
	message[TOKEN_EXPIRE_ERROR] = "token失效，请重新登陆"
	message[TOKEN_GENERATE_ERROR] = "生成token失败"
	message[DB_ERROR] = "数据库繁忙,请稍后再试"
	message[REQUEST_VERIFY_ERROR] = "请求验证错误"
	message[USER_EXIST] = "该用户已存在"
	message[USER_REGISTRT_FAIL] = "用户注册失败"
	message[USER_LOGIN_ERR] = "用户登入失败"
	message[USER_UPDATE_PWD_ERR] = "用户修改密码失败"
	message[USER_UPDATE_STATUS_ERR] = "用户更新状态失败"
	message[USER_OPERATION_ERR] = "用户操作失败"
	message[ALREADY_EXISTS] = "已存在"
	message[DATA_NOT_FIND] = "找不到数据"
}

func MapErrMsg(errcode codes.Code) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errcode codes.Code) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
