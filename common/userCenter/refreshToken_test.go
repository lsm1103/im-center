package userCenter

import "testing"

var userC = NewUserCenter(
	"http://172.16.10.184:10090",
	"Pj-Diagnosis-C",
	"3d6c838983e145cdbff40284f7d0e04b",
	10,
)

func TestRefreshToken(t *testing.T) {
	info, err := userC.RefreshToken(&RefreshTokenReq{
		Refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzY5OTczNzUsInN1YiI6IjYxMDUzNzM4ZTBkYTRmODVhOGVjY2Q1MzQzZjFkMTc0IiwiaWF0IjoxNjM2MzkyNTc1LCJ0eXBlIjoicmVmcmVzaCIsImp0aSI6InRoaXJkU2VydmVyIiwiY2wiOiJhYWFjY2MifQ.pAO-7xWxWzgGKANDpdre820HGMDDRhWzXNkqmxc7r80",
	})
	if err != nil {
		panic(err)
	}
	t.Logf("%+v, %+v",info, err)
}