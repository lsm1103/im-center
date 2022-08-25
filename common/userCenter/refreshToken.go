package userCenter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type (
	RefreshTokenReq struct {
		Refresh_token string `json:"refresh_token,omitempty"`
	}
	RefreshTokenResp struct {
		Code int64    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			CliendId        string `json:"cliend_id"`
			AccessToken     string `json:"access_token"`
			RefreshToken    string `json:"refresh_token"`
			AccessExpireIn  int64    `json:"access_expire_in"`
			RefreshExpireIn int64    `json:"refresh_expire_in"`
		} `json:"data"`
	}

)

func (u *UserCenter)RefreshToken(req *RefreshTokenReq) (*RefreshTokenResp, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(req); err != nil { return nil, err }
	fmt.Println(payload)
	refresh, err := u.req.Post(
		fmt.Sprintf("%s/api/v1/authior/refresh?clientId=%s&clientSecret=%s&grant_type=%s", u.baseUrl, u.clientId, u.clientSecret, "refresh_token" ),
		"application/json",
		payload,
	)

	if err != nil {
		return nil, err
	}
	defer refresh.Body.Close()

	var resp RefreshTokenResp
	err = json.NewDecoder(refresh.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
