package requestsAiDph

import "testing"

func TestGetAis(t *testing.T) {
	ais, err := reqAi.GetAis(&GetAisReq{
		Current: 1,
		PageSize:   10,
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v, %+v",ais, err)
}