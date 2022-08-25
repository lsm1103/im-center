package userCenter

import (
	"net/http"
	"time"
)

type UserCenter struct {
	baseUrl string
	clientId string
	clientSecret string
	req *http.Client
}

func NewUserCenter(baseUrl string, clientId string, clientSecret string, timeout int64 ) *UserCenter {
	return &UserCenter{
		baseUrl:      baseUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		req:          &http.Client{Timeout: time.Duration(timeout) * time.Second}, //默认超时时间10秒
	}
}

func (u *UserCenter) BaseUrl() string {
	return u.baseUrl
}

func (u *UserCenter) SetBaseUrl(baseUrl string) {
	u.baseUrl = baseUrl
}

func (u *UserCenter) ClientId() string {
	return u.clientId
}

func (u *UserCenter) SetClientId(clientId string) {
	u.clientId = clientId
}

func (u *UserCenter) ClientSecret() string {
	return u.clientSecret
}

func (u *UserCenter) SetClientSecret(clientSecret string) {
	u.clientSecret = clientSecret
}

func (u *UserCenter) Req() *http.Client {
	return u.req
}

func (u *UserCenter) SetReq(req *http.Client) {
	u.req = req
}
