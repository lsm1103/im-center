package userCenter

import (
	"encoding/json"
	"fmt"
	"im-center/common/xerr"
)

type (
	GetTokenByCode struct {
		Code       string `json:"code"`
		Grant_type string `json:"grant_type,options=authorization_code|implicit|password|client_credentials,default=authorization_code"`
		Three_type string `json:"three_type,options=1|2,omitempty"`	//1: 钉钉; 2: 微信
	}
	GetTokenByCodeResp struct {
		Code int64    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			CliendId        string `json:"cliend_id"`
			AccessToken     string `json:"access_token"`
			RefreshToken    string `json:"refresh_token"`
			AccessExpireIn  int64    `json:"access_expire_in"`
			RefreshExpireIn int64    `json:"refresh_expire_in"`
			TokenType       string `json:"token_type"`
			Info            struct {
				CreateTime  string `json:"createTime"`
				UpdateTime  string `json:"updateTime"`
				Id          string `json:"id"`
				RealName    string `json:"realName"`
				Gender      string `json:"gender"`
				//Password    string `json:"password"`
				Superuser   int64    `json:"superuser"`
				PhoneNumber string `json:"phoneNumber"`
				Flag        int64    `json:"flag"`
				UserInfo    struct {
					CreateTime   string      `json:"createTime"`
					UpdateTime   string      `json:"updateTime"`
					Id           int64         `json:"id"`
					UserId       string      `json:"userId"`
					Server       string      `json:"server"`
					NickName     string `json:"nick_name"`
					UserEmail    string `json:"userEmail"`
					OtherContact string `json:"otherContact"`
					ExpertImage  string `json:"expertImage"`
					Hospital     string      `json:"hospital"`
					Department   string      `json:"department"`
					JobNumber    string      `json:"jobNumber"`
					OwnerUser    string `json:"ownerUser"`
				} `json:"user_info"`
			} `json:"info"`
		} `json:"data"`
	}

	UserInfo struct {
		CliendId        string `json:"cliend_id"`
		AccessToken     string `json:"access_token"`
		RefreshToken    string `json:"refresh_token"`
		AccessExpireIn  int64  `json:"access_expire_in"`
		RefreshExpireIn int64  `json:"refresh_expire_in"`
		TokenType       string `json:"token_type"`
		Info            struct {
			CreateTime  string `json:"createTime"`
			UpdateTime  string `json:"updateTime"`
			Id          string `json:"id"`
			RealName    string `json:"realName"`
			Gender      string `json:"gender"`
			Superuser   int64  `json:"superuser"`
			PhoneNumber string `json:"phoneNumber"`
			Flag        int64  `json:"flag"`
			UserInfo    struct {
				CreateTime   string `json:"createTime"`
				UpdateTime   string `json:"updateTime"`
				Id           int64  `json:"id"`
				UserId       string `json:"userId"`
				Server       string `json:"server"`
				NickName     string `json:"nick_name"`
				UserEmail    string `json:"userEmail"`
				OtherContact string `json:"otherContact"`
				ExpertImage  string `json:"expertImage"`
				Hospital     string `json:"hospital"`
				Department   string `json:"department"`
				JobNumber    string `json:"jobNumber"`
				OwnerUser    string `json:"ownerUser"`
			} `json:"user_info"`
		} `json:"info"`
	}
)

func (u *UserCenter)GetTokenByCode(req *GetTokenByCode) (*GetTokenByCodeResp, error) {
	token, err := u.req.Get(fmt.Sprintf("%s/api/v1/authior/token?code=%s&grant_type=%s&clientId=%s&clientSecret=%s",
		u.baseUrl, req.Code, req.Grant_type, u.clientId, u.clientSecret) )
	if err != nil { return nil, err }
	defer token.Body.Close()

	var resp GetTokenByCodeResp
	err = json.NewDecoder(token.Body).Decode(&resp)
	//fmt.Printf("resp:%+v",resp)
	if err != nil { return nil, err }
	if resp.Code != int64(200) && resp.Code != int64(201) {
		return nil, xerr.NewErrMsg(resp.Msg)
	}
	return &resp, nil
}

// 用来实现拿第三方的code去换用户中心的code，然后拿用户中心code换token
//type GetTokenByThreeResp struct {
//	Code int64    `json:"code"`
//	Msg  string `json:"msg"`
//	Data string `json:"data"`
//}
//var token *http.Response
//var err error
//if req.Three_type == "" {
//	token, err = u.req.Get(fmt.Sprintf("%s/api/v1/authior/token?code=%s&grant_type=%s&clientId=%s&clientSecret=%s",
//		u.baseUrl, req.Code, req.Grant_type, u.clientId, u.clientSecret) )
//} else {
//	//先用第三方code获取用户中心code
//	body := new(bytes.Buffer)
//	if err = json.NewEncoder(body).Encode(map[string]string{
//		"code":req.Code,
//		"state":u.clientId,
//		"type": req.Three_type,
//	}); err != nil { return nil, err }
//	reqt, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dingding", u.baseUrl), body)
//	if err != nil { return nil, err }
//	reqt.Header.Add("Cookies", fmt.Sprintf("clientId=%s",u.clientId))
//	code, err := u.req.Do(reqt)
//	if err != nil { return nil, err }
//	var resp_ GetTokenByThreeResp
//	err = json.NewDecoder(code.Body).Decode(&resp_)
//	if err != nil { return nil, err }
//	if resp_.Code != int64(200) && resp_.Code != int64(201) {
//		fmt.Printf("resp_:%+v",resp_)
//		return nil, xerr.NewErrMsg(resp_.Msg)
//	}
//	//获取token
//	token, err = u.req.Get(fmt.Sprintf("%s/api/v1/authior/token?code=%s&grant_type=%s&clientId=%s&clientSecret=%s",
//		u.baseUrl, resp_.Data, req.Grant_type, u.clientId, u.clientSecret) )
//}
