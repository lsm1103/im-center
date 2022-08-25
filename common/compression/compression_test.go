package compression

import "testing"

func TestCompression(t *testing.T) {
	err := Compression("/Users/xm/Desktop/work_project/im-center/dataFile/sliceFile/2022/06/08/a37720563f0c0761da83e678af44668a", "/Users/xm/Desktop/11/a37720563f0c0761da83e678af44668a.zip")
	//err := Compression("/Users/xm/Desktop/11/dd", "/Users/xm/Desktop/11/dd.zip")
	if err != nil {
		panic(err)
	}
	t.Logf("err:%+v",err)
}
