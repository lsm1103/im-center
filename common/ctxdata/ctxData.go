package ctxdata

import (
	"context"
	"google.golang.org/grpc/metadata"
	"im-center/common/globalkey"
	"strconv"
)

//时间模版 Y-Y-Y H:M:S
var TimeFmt string = "2006-01-02 15:04:05"

// 从ctx获取值
var CtxKeyUserData string = "UserData"
var CtxKeyGrantType string = "GrantType"

var CtxKeyUserId string = "UserId"
var CtxKeyRegisterDevice string = "RegisterDevice"
var CtxKeyIdentityType string = "IdentityType"
var CtxKeyIdentifier string = "Identifier"
var CtxKeyNickname string = "Nickname"
var CtxKeySex string = "Sex"
var CtxKeyIco string = "Ico"
var CtxKeyStatus string = "Status"

//机器类型
var MachineTypes = map[string]int64{
	"Other": 0, "Android": 1, "IOS": 2, "Windows": 3, "Mac OS X": 4, "Web": 5,
}

//需验证的Api列表
var VerifyTypes = map[string][]string{
	"DeviceApis":    []string{"/dddddd"},
	"SecretKeyApis": []string{"/dddddd"},
	"TokenApis":     []string{
		"/doc",
		"/docData",
		"/user_center/v1/userRegister",
		"/user_center/v1/userLogin",
		"/file_center/v0/*",
		"/predict/v1/*",
		"/auth/v1/getToken",
		"/auth/v1/getAuthInfo",
	},
}

//状态类型
var StatusTypes = map[string]int64{
	"del":globalkey.Del,
	"freeze":globalkey.UserFreeze,
}

func GetRegisterDeviceFromCtx(ctx context.Context) int64 {
	if md,ok := metadata.FromIncomingContext(ctx); ok{
		val := md.Get(CtxKeyRegisterDevice)
		if val == nil || len(val) < 0{
			return 0
		}
		r,err := strconv.ParseInt(val[0], 10, 64)
		if err != nil{
			return 0
		}
		return r
	}
	return 0
}

func GetKeyFromCtx(ctx context.Context, key string) string {
	if md,ok := metadata.FromIncomingContext(ctx); ok{
		val := md.Get(key)
		if val == nil || len(val) < 0{
			return ""
		}
		return val[0]
	}
	return ""
}

func GetUserIdFromCtx(ctx context.Context) int64 {
	if md,ok := metadata.FromIncomingContext(ctx); ok{
		val := md.Get(CtxKeyUserId)
		if val == nil || len(val) < 0{
			return 0
		}
		r,err := strconv.ParseInt(val[0], 10, 64)
		if err != nil{
			return 0
		}
		return r
	}
	return 0
}
