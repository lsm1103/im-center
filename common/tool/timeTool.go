package tool

import "time"

var TimeFmt string = "2006-01-02 15:04:05"
var TimeSlimFmt string = "20060102-150405"
var Loc, _ = time.LoadLocation("Asia/Shanghai") //上海

func FmtTime(t time.Time) string {
	return t.In(Loc).Format(TimeFmt)
}

func StrToTime(str string) time.Time {
	location, err := time.ParseInLocation(TimeFmt, str, Loc)
	if err != nil {
		return time.Time{}
	}
	return location
}

func GetNowTime(timeFmt string) string {
	return time.Now().In(Loc).Format(timeFmt)
}

func GetDefNowTime() string {
	return time.Now().In(Loc).Format(TimeFmt)
}

func GetUintNowTime() uint64 {
	return uint64(time.Now().In(Loc).Unix())
}