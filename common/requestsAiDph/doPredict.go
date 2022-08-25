package requestsAiDph

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type (
	DoPredictReq struct {
		Curd struct {
			DataUid       string `json:"data_uid"`
			AiType        string `json:"ai_type"`
			//UseUid        string `json:"use_uid"`
			SourcePath    string `json:"source_path"`
			//PredictPath   string `json:"predict_path"`
			//MaskPath      string `json:"mask_path"`
			//PredictStatus int64    `json:"predict_status"`
			//PredictInfo   string `json:"predict_info"`
			//Kv            string `json:"kv"`
			//Status int64 `json:"status"`
		} `json:"curd"`
		Priority      int64    `json:"priority"`
		//CallbackUrl   string `json:"callback_url"`
		//CallbackToken string `json:"callback_token"`
	}
	DoPredictResp struct {
		Code int64 `json:"code"`
		Data struct {
			AiType     string `json:"ai_type"`
			CreateTime string `json:"create_time"`
			DataUid    int64    `json:"data_uid"`
			Id         int64    `json:"id"`
			Kv         struct {
			} `json:"kv"`
			MaskPath      string `json:"mask_path"`
			PredictInfo   string `json:"predict_info"`
			PredictPath   string `json:"predict_path"`
			PredictStatus int64    `json:"predict_status"`
			SourcePath    int64    `json:"source_path"`
			Status        int64    `json:"status"`
			TaskId        string `json:"task_id"`
			Uid           string `json:"uid"`
			UpdateTime    string `json:"update_time"`
			UseUid        string `json:"use_uid"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
)

func (r *ReqAiDph) DoPredict(req *DoPredictReq) (*DoPredictResp, error) {
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(req); err != nil { return nil, err }
	predict, err := http.NewRequest("POST", r.base_url + "/v0/predict/predict", reqBody)
	predict.Header.Set("Content-Type", "application/json")
	predict.Header.Set("app-key", r.uid)
	predict.Header.Set("verifyType", "NoVerify")
	predict.Header.Set("sign", "")
	resp, err := r.req.Do(predict)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bd DoPredictResp
	err = json.NewDecoder(resp.Body).Decode(&bd)
	return &bd, nil
}


//{"code":100200,
//"data":{
//"ai_type":"213062728c8811ebb4b5e454e8c157e1",
//"create_time":"2022-05-11 19:21:22",
//"data_uid":122,
//"id":770,
//"kv":{},
//"mask_path":"",
//"predict_info":"",
//"predict_path":"",
//"predict_status":0,
//"source_path":1,
//"status":0,
//"task_id":"35c79364-53c0-465c-9aa6-ab9c41a2a5f8",
//"uid":"7c4f88cad11c11ec83ce024271000d3e",
//"update_time":"2022-05-11 19:21:22",
//"use_uid":"test_local"
//},
//"msg":"ok"
//}
