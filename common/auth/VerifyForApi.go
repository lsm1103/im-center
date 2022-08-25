package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"im-center/common/ctxdata"
	"im-center/common/xerr"
)

//是否需要验证Api
func IsVerifyForApi(api string, verifyType string, isReverse bool) bool {
	exist := false
	api_list := strings.Split(api,"/")
	for _, item := range ctxdata.VerifyTypes[verifyType] {
		if item == api {
			exist = true
			break
		}
		item_list := strings.Split(item,"/")
		if item_list[len(item_list)-1] == "*" {
			api_ := strings.Join(api_list[:len(item_list)-1],"/")
			item_ := strings.Join(item_list[:len(item_list)-1],"/")
			if api_ == item_ {
				exist = true
				break
			}
		}
	}
	if isReverse {
		exist = ! exist
	}
	return exist
}

//设备验证
func VerifyDevice(header *http.Header) error {
	// 验证设备ip/mac,实现一些白名单、黑名单的需求
	ip := header.Get("ip")
	mac := header.Get("mac")
	fmt.Printf("ip:%v, mac:%v\n", ip, mac)
	return nil
	//return errors.Wrapf(xerr.NewErrCode(xerr.TOKEN_EXPIRE_ERROR), "VerifySecretKey 该设备不合法") //该设备不合法
}

//密钥验证，特定接口
func VerifySecretKey(header *http.Header) error {
	ClientID := header.Get("ClientID")
	Signature := header.Get("Signature")
	fmt.Printf("ClientID:%v, Signature:%v\n", ClientID, Signature)
	return nil
}

//令牌验证，todo 现做AccessSecret固定，后面可以考虑每个人的AccessSecret不一样
func VerifyToken(r *http.Request, accessSecret string) (map[string]string, error) {
	UserData,err := ParseVerifyJwtToken(r,accessSecret)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.TOKEN_EXPIRE_ERROR), "%v", err)
	}
	return UserData, err
}

/*
TODO
	1、设备信息走deviceRegister接口发到用户服务（客户端打开页面的第一时间(最好保证在用户注册之前拿到设备注册返回的设备id)、
	2、每隔一定时间(这个时候也可以尝试去查询该用户有没有绑定该设备，做补偿) ）
	3、普通的接口请求只需要验证ua、ip等，所以令牌和设备验证做在同一个接口
*/
