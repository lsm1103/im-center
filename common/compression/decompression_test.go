package compression

import "testing"

func TestDeCompression(t *testing.T) {
	err := DeCompression(
		"/Users/xm/Desktop/work_project/im-center/dataFile/normalFile/0e3b53a70244d8fa1b6426335f3ed700.zip",
		//"/Users/xm/Desktop/work_project/im-center/dataFile/normalFile/5d5ffeb51ca3f79595a77f6cb71494b7.zip",
		"/Users/xm/Desktop/work_project/im-center/dataFile/sliceFile/2022/06/13/0e3b53a70244d8fa1b6426335f3ed700",
		false,
		true)
	//err := DeCompression("/Users/xm/Desktop/11/dd-.zip", "/Users/xm/Desktop/11/dd", false, true)
	//err := DeCompression("/Users/xm/Desktop/test/徐思源_20220513090732628.zip", "/Users/xm/Desktop/11/dd", true, true)
	if err != nil {
		panic(err)
	}
	t.Logf("err:%+v",err)
}
