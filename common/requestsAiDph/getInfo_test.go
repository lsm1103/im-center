package requestsAiDph

import (
	"encoding/json"
	"testing"
)

var reqAi = NewReqAiDph(
	"本地测试",
	ReqAiDphCfg{
		Uid:      "test_local",
		Base_url: "http://172.16.10.91:8802",
		Ais:      "ais",
		Timeout:  1000,
	},
	nil,
)

func TestGetInfo(t *testing.T) {
	info, err := reqAi.GetInfo(&GetInfoReq{
		//TimeRange:  "1949-01-01 00:00:00/2022-05-06 18:12:21",
		//In:         map[string]interface{}{},
		QueryType:  "and",
		QueryFuzzy: false,
		Curd: map[string]interface{}{
			"data_uid":616,
		},
	})
	if err != nil {
		panic(err)
	}
	t.Logf("info:%+v, %+v\n",info, err)
	var q PredictKV
	err = json.Unmarshal([]byte(info.Kv), &q)
	t.Logf("q:%+v, %+v\n",q, err)
	t.Logf("q:%+v, %+v\n",q.EstimateTotal, q.EstimateTime)

}


//type PredictInfo struct {
//	DataUid    string `json:"data_uid"`
//	AiType     string `json:"ai_type"`
//	UseUid     string `json:"use_uid"`
//	SourcePath string `json:"source_path"`
//	Priority                int64     `json:"priority"`
//	TaskId                  string  `json:"task_id"`
//	VerifyType              string  `json:"verifyType"`
//	DataPreprocessingSTime  string  `json:"DataPreprocessingSTime"`
//	DataPreprocessingETime  string  `json:"DataPreprocessingETime"`
//	PredictStatus           int64     `json:"predict_status"`
//	FilterPath              string  `json:"filter_path"`
//	PredictPath             string  `json:"predict_path"`
//	FattyPath               string  `json:"fatty_path"`
//	StmId                   string  `json:"stmId"`
//	EstimateTime            float64 `json:"estimate_time"`
//	EstimateTotal           float64 `json:"estimate_total"`
//	PredictETime            string  `json:"PredictETime"`
//	MaskPath                string  `json:"mask_path"`
//	DataPostProcessingSTime string  `json:"DataPostProcessingSTime"`
//	PredictInfo             string  `json:"predict_info"`
//	DataPostProcessingETime string  `json:"DataPostProcessingETime"`
//	Mask2Report             string  `json:"mask2report"`
//	Mask2Seg                struct {
//		Res     int64      `json:"res"`
//		SegPath []string `json:"seg_path"`
//	} `json:"mask2seg"`
//	Mask2Json struct {
//		Len      int64    `json:"len"`
//		JsonPath string `json:"json_path"`
//	} `json:"mask2json"`
//}

//info:
//{
//	CreateTime:2022-05-06 18:06:27
//	UpdateTime:2022-05-06 18:07:34
//	Id:681
//	Uid:31193b90cd2411ecaa62024271000d3e
//	DataUid:486
//	UseUid:diagnose
//	PredictStatus:6
//	PredictInfo:数据后处理完成
//	AiType:213062728c8811ebb4b5e454e8c157e1
//	SourcePath:diagnose/pjkj_2021_1203/studies/2022/04/20/8d4ac332903a4332826a55a661f5904d_study
//	PredictPath:three_series/ct/2022/04/20/ct_31193b90cd2411ecaa62024271000d3e_three.npz
//	MaskPath:three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three.npz
//	Kv:{
//		"mask2report": "",
//		"mask2seg": {
//			"res": 1,
//			"seg_path": [
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_2_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_2_liver_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_6_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_6_liver_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_9_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_9_liver_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_12_seg.dcm",
//			"three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_12_liver_seg.dcm"
//			]
//		},
//		"mask2json": {
//			"len": 13599,
//			"json_path": "three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three.json"
//		},
//		"DataPostProcessingETime": "2022-05-06 18:07:34",
//		"PredictETime": "2022-05-06 18:07:32",
//		"DataPostProcessingSTime": "2022-05-06 18:07:32",
//		"estimate_time": 158.0,
//		"estimate_total": 158.0,
//		"stmId": "1651831590528-0",
//		"priority": 1,
//		"DataPreprocessingETime": "2022-05-06 18:06:29",
//		"filter_path": "dicom_filter/diagnose/pjkj_2021_1203/studies/2022/04/20/8d4ac332903a4332826a55a661f5904d_study_threeFilter",
//		"fatty_path": ""
//	}
//	Status:1
//}
////q:
//{
//	DataUid: AiType: UseUid: SourcePath:
//	Priority:1
//	TaskId: VerifyType: DataPreprocessingSTime:
//	DataPreprocessingETime:2022-05-06 18:06:29
//	PredictStatus:0
//	FilterPath:dicom_filter/diagnose/pjkj_2021_1203/studies/2022/04/20/8d4ac332903a4332826a55a661f5904d_study_threeFilter
//	PredictPath:
//	FattyPath:
//	StmId:1651831590528-0
//	EstimateTime:158
//	EstimateTotal:158
//	PredictETime:2022-05-06 18:07:32
//	MaskPath:
//	DataPostProcessingSTime:2022-05-06 18:07:32
//	PredictInfo: DataPostProcessingETime:2022-05-06 18:07:34
//	Mask2Report:
//	Mask2Seg:{
//		Res:1
//		SegPath:[
//		three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_2_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_2_liver_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_6_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_6_liver_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_9_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_9_liver_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_12_seg.dcm three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three_12_liver_seg.dcm]}
//	Mask2Json:{
//		Len:13599
//		JsonPath:three_series/mask/2022/04/20/31193b90cd2411ecaa62024271000d3e_three.json
//	}
//}
