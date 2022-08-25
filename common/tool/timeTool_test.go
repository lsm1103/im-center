package tool

import "testing"

func TestGetNowTime(t *testing.T) {
	s := GetNowTime(TimeSlimFmt)
	t.Log(s)
}