package requestsAiDph

import "testing"

func TestDoPredict(t *testing.T) {
	info, err := reqAi.DoPredict(&DoPredictReq{
		Curd: struct {
			DataUid       string `json:"data_uid"`
			AiType        string `json:"ai_type"`
			SourcePath    string `json:"source_path"`
		}{
			DataUid: "122",
			AiType: "213062728c8811ebb4b5e454e8c157e1",
			SourcePath: "1",
		},
		Priority:      1,
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v, %+v",info, err)
}