package userCenter

import "testing"

func TestVerifyToken(t *testing.T) {
	info, err := userC.VerifyToken(&VerifyTokenReq{
		ClientId: "clientId=Pj-Diagnosis-C",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2MjY5OTcsInN1YiI6ImRjZGVkMGE2ZGUxODQwOGI5MDJjM2FjYmFlYmM4MTJlIiwiaWF0IjoxNjUyMzY3Nzk3LCJ0eXBlIjoiYWNjZXNzIiwianRpIjoidXNlciIsImNsIjoiUGotRGlhZ25vc2lzLUMifQ.XsEKQaKYFaZrOTQoRFdNvWOHehxD82n9LadHELQrNZg",
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v, %+v",info, err)
}