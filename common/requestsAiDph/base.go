package requestsAiDph

import (
	"net/http"
	"im-center/common/cacheHandle"
	"time"
)

type ReqAiDph struct {
	uid                 string
	name                string
	base_url			string
	account             string
	passwd              string
	ais                 string
	secret_key          string
	req					*http.Client
	CacheHandle *cacheHandle.CecheHandle
}

type ReqAiDphCfg struct {
	Uid        string
	Base_url   string
	Account    string
	Passwd     string
	Ais        string
	Secret_key string
	Timeout    int64
}

func NewReqAiDph(name string, cfg ReqAiDphCfg, cacheHandle *cacheHandle.CecheHandle ) *ReqAiDph {
	return &ReqAiDph{
		name: name,
		uid: cfg.Uid,
		base_url: cfg.Base_url,
		account: cfg.Account,
		passwd: cfg.Passwd,
		ais: cfg.Ais,
		secret_key: cfg.Secret_key,
		req: &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second }, //默认超时时间10秒
		CacheHandle: cacheHandle,
	}
}

func (r *ReqAiDph) Base_url() string {
	return r.base_url
}

func (r *ReqAiDph) SetBase_url(base_url string) {
	r.base_url = base_url
}

func (r *ReqAiDph) Name() string {
	return r.name
}

func (r *ReqAiDph) SetName(name string) {
	r.name = name
}

func (r *ReqAiDph) Uid() string {
	return r.uid
}

func (r *ReqAiDph) SetUid(uid string) {
	r.uid = uid
}

func (r *ReqAiDph) Account() string {
	return r.account
}

func (r *ReqAiDph) SetAccount(account string) {
	r.account = account
}

func (r *ReqAiDph) Passwd() string {
	return r.passwd
}

func (r *ReqAiDph) SetPasswd(passwd string) {
	r.passwd = passwd
}

func (r *ReqAiDph) Ais() string {
	return r.ais
}

func (r *ReqAiDph) SetAis(ais string) {
	r.ais = ais
}

func (r *ReqAiDph) Secret_key() string {
	return r.secret_key
}

func (r *ReqAiDph) SetSecret_key(secret_key string) {
	r.secret_key = secret_key
}
