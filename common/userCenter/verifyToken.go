package userCenter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"im-center/common/xerr"
)

type (
	VerifyTokenReq struct {
		ClientId string `json:"client_id,omitempty"`
		Token    string `json:"token,omitempty"`
	}
	VerifyTokenResp struct {
		Code int64    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			CreateTime  string `json:"createTime"`
			UpdateTime  string `json:"updateTime"`
			Id          string `json:"id"`
			RealName    string `json:"realName"`
			Gender      string `json:"gender"`
			//Password    string `json:"password"`
			Superuser   int64  `json:"superuser"`
			PhoneNumber string `json:"phoneNumber"`
			Flag        int64  `json:"flag"`
			Info        struct {
				CreateTime   string      `json:"createTime"`
				UpdateTime   string      `json:"updateTime"`
				Id           int64       `json:"id"`
				UserId       string      `json:"userId"`
				Server       string      `json:"server"`
				NickName     string 	 `json:"nick_name"`
				UserEmail    string 	 `json:"userEmail"`
				OtherContact string 	 `json:"otherContact"`
				ExpertImage  string 	 `json:"expertImage"`
				Hospital     string      `json:"hospital"`
				Department   string      `json:"department"`
				JobNumber    string      `json:"jobNumber"`
				OwnerUser    string 	 `json:"ownerUser"`
			} `json:"info"`
		} `json:"data"`
	}
)

func (u *UserCenter)VerifyToken(req *VerifyTokenReq) (*VerifyTokenResp, error) {
	reqt, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/authior/verify/token", u.baseUrl), nil)
	reqt.Header.Set("Cookies", req.ClientId)
	reqt.Header.Set("token", req.Token)
	token, err := u.req.Do(reqt)
	if err != nil { return nil, err }
	defer token.Body.Close()

	var resp VerifyTokenResp
	err = json.NewDecoder(token.Body).Decode(&resp)
	fmt.Printf("resp:%+v",resp)
	if err != nil { return nil, err }
	if resp.Code != int64(200) && resp.Code != int64(201) {
		return nil, xerr.NewErrMsg(resp.Msg)
	}
	return &resp, nil
}

//获取缓存
//get, err := l.svcCtx.RedisClient.Get(globalkey.CacheUserInfoByTokenKey)
////fmt.Printf("get:%+v, err:%+v\n",get, err )
//var q userCenter.UserInfo
//err = json.Unmarshal([]byte(get), &q)
////fmt.Printf("q:%+v, err:%+v\n",q, err )