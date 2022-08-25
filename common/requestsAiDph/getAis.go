package requestsAiDph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"im-center/common/ctxdata"
	"im-center/common/globalkey"
	"im-center/common/tool"
	"time"
)

type GetAisReq struct {
	TimeRange  string `json:"time_range"`
	Current    int64  `json:"current"`
	PageSize   int64  `json:"page_size"`
}

type AiItem struct {
	CreateTime         string `json:"create_time"`
	UpdateTime         string `json:"update_time"`
	Id                 int64    `json:"id"`
	Uid                string `json:"uid"`
	Name               string `json:"name"`
	Version            string `json:"version"`
	Url                string `json:"url"`
	Account            string `json:"account"`
	Password           string `json:"password"`
	DataPreprocessing  string `json:"data_preprocessing"`
	DataPostProcessing string `json:"data_post_processing"`
	Kv                 string `json:"kv"`
	Status             int64    `json:"status"`
}

type GetAis struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []AiItem `json:"list"`
		Pagination struct {
			Start      int64 `json:"start"`
			Current    int64 `json:"current"`
			TotalPages int64 `json:"total_pages"`
			PageSize   int64 `json:"page_size"`
			Total      int64 `json:"total"`
		} `json:"pagination"`
	} `json:"data"`
}

type GetAisResp struct {
	AiMap map[string]AiItem `json:"list"`
	Pagination struct {
		Start      int64 `json:"start"`
		Current    int64 `json:"current"`
		TotalPages int64 `json:"total_pages"`
		PageSize   int64 `json:"page_size"`
		Total      int64 `json:"total"`
	} `json:"pagination"`
}

func (r *ReqAiDph) GetAis(req *GetAisReq) (*GetAisResp, error) {
	req.TimeRange = fmt.Sprintf("1949-01-01 00:00:00/%s", time.Now().In(tool.Loc).Format(ctxdata.TimeFmt))
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(req); err != nil { return nil, err }
	resp, err := r.req.Post(r.base_url + "/v0/Ai/gets", "application/json", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var getAis GetAis
	if err = json.NewDecoder(resp.Body).Decode(&getAis); err != nil { return nil, err }
	getAisResp := GetAisResp{ Pagination: getAis.Data.Pagination, AiMap: map[string]AiItem{}}
	for _,item := range getAis.Data.List {
		getAisResp.AiMap[item.Name] = item
		//getAisResp.AiMap[item.Uid] = item
	}
	return &getAisResp, nil
}

func (r *ReqAiDph) GetAisCeche(req *GetAisReq) (*GetAisResp, error) {
	var resp GetAisResp
	//1、获取算法权限列表：先查询缓存，如果命中就返回；未命中就查询算法调度服务并更新缓存后返回
	if err := r.CacheHandle.Get(&resp, globalkey.CacheAiInfoKey, func(val interface{}) error {
		info, err := r.GetAis(req)
		if err != nil { return err }
		val.(*GetAisResp).AiMap = info.AiMap
		val.(*GetAisResp).Pagination = info.Pagination
		return nil
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}