package userCenter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"im-center/common/xerr"
)

type (
	GetUserInfoReq struct {
		UserId        string `json:"user_id"`
		Authorization string `json:"authorization"`
	}
)

func (u *UserCenter)GetUserInfo(in *GetUserInfoReq) (*GetTokenByCodeResp, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/%s", u.baseUrl, in.UserId), nil)
	if err != nil { return nil, err }
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s",in.Authorization))
	req.Header.Add("Cookies", fmt.Sprintf("clientId=%s",u.clientId))
	token, err := u.req.Do(req)
	if err != nil { return nil, err }
	var resp GetTokenByCodeResp
	if err = json.NewDecoder(token.Body).Decode(&resp); err != nil { return nil, err }
	fmt.Printf("resp:%+v",resp)
	if resp.Code != int64(200) && resp.Code != int64(201) {
		return nil, xerr.NewErrMsg(resp.Msg)
	}
	return &resp, nil
}
