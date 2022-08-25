package requestsAiDph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"im-center/common/ctxdata"
	"im-center/common/tool"
	"time"
)

type (
	GetInfoReq struct {
		TimeRange string `json:"time_range,omitempty"`
		In  map[string]interface{} `json:"in_"`
		//{
		//	Name string        `json:"name,omitempty"`
		//	List []interface{} `json:"list_,omitempty"`
		//}
		QueryType  string `json:"query_type"`
		QueryFuzzy bool   `json:"query_fuzzy"`
		Curd  map[string]interface{} `json:"curd"`
		//{
		//	Data_uid string `json:"data_uid,omitempty"`
		//	Ai_type     string `json:"ai_type,omitempty"`
		//	Use_uid      string `json:"use_uid,omitempty"`
		//	Source_path  string `json:"source_path,omitempty"`
		//	Predict_path   string `json:"predict_path,omitempty"`
		//	Mask_path      string `json:"mask_path,omitempty"`
		//	Predict_status int64  `json:"predict_status,omitempty"`
		//	Predict_info string `json:"predict_info,omitempty"`
		//	Kv           string `json:"kv,omitempty"`
		//	Status int64  `json:"status,omitempty"`
		//	Uid    string `json:"uid,omitempty"`
		//}
	}
	GetInfoResp struct {
		Code int64    `json:"code"`
		Msg  string `json:"msg"`
		//Data map[string]interface{}  `json:"data"`
		Data Info
	}
	Info struct {
		CreateTime    string `json:"create_time"`
		UpdateTime    string `json:"update_time"`
		Id            int64    `json:"id"`
		Uid           string `json:"uid"`
		DataUid       string `json:"data_uid"`
		UseUid        string `json:"use_uid"`
		PredictStatus int64    `json:"predict_status"`
		PredictInfo   string `json:"predict_info"`
		AiType        string `json:"ai_type"`
		SourcePath    string `json:"source_path"`
		PredictPath   string `json:"predict_path"`
		MaskPath      string `json:"mask_path"`
		Kv            string `json:"kv"`
		Status        int64    `json:"status"`
	}
	PredictKV struct {
		Mask2Report string `json:"mask2report"`
		Mask2Seg    struct {
		Res     int      `json:"res"`
		SegPath []string `json:"seg_path"`
		} `json:"mask2seg"`
		Mask2Json struct {
		Len      int    `json:"len"`
		JsonPath string `json:"json_path"`
		} `json:"mask2json"`
		DataPostProcessingETime string  `json:"DataPostProcessingETime"`
		PredictETime            string  `json:"PredictETime"`
		DataPostProcessingSTime string  `json:"DataPostProcessingSTime"`
		EstimateTime            float64 `json:"estimate_time"`
		EstimateTotal           float64 `json:"estimate_total"`
		StmId                   string  `json:"stmId"`
		Priority                int     `json:"priority"`
		DataPreprocessingETime  string  `json:"DataPreprocessingETime"`
		FilterPath              string  `json:"filter_path"`
		FattyPath               string  `json:"fatty_path"`
	}
)

func (r *ReqAiDph) GetInfo(req *GetInfoReq) (*Info, error) {
	req.TimeRange = fmt.Sprintf("1949-01-01 00:00:00/%s", time.Now().In(tool.Loc).Format(ctxdata.TimeFmt))
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(req); err != nil { return nil, err }
	resp, err := r.req.Post(r.base_url + "/v0/Predict/get", "application/json", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var getInfoResp GetInfoResp
	if err = json.NewDecoder(resp.Body).Decode(&getInfoResp); err != nil { return nil, err }
	return &getInfoResp.Data, nil
}

//reqBody, err := json.Marshal(req)
//if err != nil {
//	return nil, err
//}
//payload := strings.NewReader(string(reqBody))
//payload := strings.NewReader(`{
//  "in_": {},
//  "query_type": "and",
//  "query_fuzzy": false,
//  "curd": {}
//}`)

//client := &http.Client{}
//req1, err := http.NewRequest("POST", r.base_url + "/v0/Predict/get", payload)
//resp, err := client.Do(req1)