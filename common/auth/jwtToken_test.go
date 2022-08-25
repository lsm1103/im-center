package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestGenerateJwtToken(t *testing.T) {
	now := time.Now().Unix()
	userData := map[string]string{"Userid": "1111", "UserName": "2222"}
	s, err := GenerateJwtToken("wehovf22", now, 100*60, userData)
	t.Log(s, err)
}

/*
=== RUN   TestGenerateJwtToken
    jwtToken_test.go:13: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGF0YSI6eyJVc2VyTmFtZSI6IjIyMjIiLCJVc2VyVWlkIjoiMTExMSJ9LCJleHAiOjE2NDY2NDk4MzksImlhdCI6MTY0NjY0MzgzOX0.QIqkFzeOcOchFNKWjDw3VLnQp13WJz9owuXSESsB_ak <nil>
--- PASS: TestGenerateJwtToken (0.00s)
PASS
*/

func TestVerifyJwtToken(t *testing.T) {
	r := &http.Request{}
	r.Header = http.Header{
		"Authorization":   {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGF0YSI6eyJVc2VyTmFtZSI6IjIyMjIiLCJVc2VyVWlkIjoiMTExMSJ9LCJleHAiOjE2NDY2NDk4MzksImlhdCI6MTY0NjY0MzgzOX0.QIqkFzeOcOchFNKWjDw3VLnQp13WJz9owuXSESsB_ak"},
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"en-us"},
		"Foo":             {"Bar", "two"},
	}
	s, err := ParseVerifyJwtToken(r, "")
	t.Log(s, err)
}
