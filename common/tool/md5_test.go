package tool

import (
	//"FtpClient/common"
	//"bytes"
	//"fmt"
	//"net/http"
	"testing"
)

func TestMd5(t *testing.T) {
	//uploadDir := "/Users/xm/Desktop/11/dd-.zip"
	uploadDir := "/Users/xm/Desktop//3c3341d294e2d628572096e381853dc3.zip"
	//3c3341d294e2d628572096e381853dc3
	//uploadDir := "/Users/xm/Desktop/11/md5-test/37268a65d67dd336c2650a8e6fa36350.zip"
	//md5: 161632b7e521a270d5b1bb6ed6d98e85
	s, err := md5Sum(uploadDir, false)
	t.Log(s, err)
}

func TestMd5_dir(t *testing.T) {
	uploadDir := "/Users/xm/Desktop/work_project/im-center/dataFile/sliceFile/2022/06/13/3c3341d294e2d628572096e381853dc3"
	//5a0bcaa8dcde35f383d7342a4889afea
	//uploadDir := "/Users/xm/Desktop/work_project/im-center/dataFile/sliceFile/2022/06/13/37268a65d67dd336c2650a8e6fa36350"
	//md5: 37268a65d67dd336c2650a8e6fa36350
	//161632b7e521a270d5b1bb6ed6d98e85
	s, err := md5Sum(uploadDir, true)
	t.Log(s, err)
}

func TestMd5Slice(t *testing.T) {
	uploadPath := "/Users/xm/Desktop/11/md5-test/37268a65d67dd336c2650a8e6fa36350.zip"
	//161632b7e521a270d5b1bb6ed6d98e85
	s, err := md5SumSlice(uploadPath)
	t.Log(s, err)
}

func TestMergeSlice(t *testing.T) {
	uploadPath := "/Users/xm/Desktop/work_project/im-center/dataFile/sliceFile/2022/06/13/5d5ffeb51ca3f79595a77f6cb71494b7"
	mergeSlice(uploadPath)
}

//func TestUpSliceFile(t *testing.T)  {
//	for i := startIndex; i < u.SliceNum; i++ {
//		tmpData := make([]byte, u.SliceBytes)
//		nr, err := fh.Read(tmpData[:])
//		if err != nil {
//			fmt.Printf("read file error\n")
//			return err
//		}
//		if md5sum == "" {
//			hash.Write(tmpData[:nr])
//		}
//		tmpData = tmpData[:nr]
//	}
//	targetUrl := common.BaseUrl + "uploadBySlice"
//	//fmt.Printf("fid: %s, index: %d\n", part.Fid, part.Index)
//
//	reqBody := new(bytes.Buffer)
//	json.NewEncoder(reqBody).Encode(part)
//
//	req, err := http.NewRequest("POST", targetUrl, reqBody)
//	req.Header.Set("Content-Type", "application/json")
//
//	resp, err := (&http.Client{}).Do(req)
//	if err != nil {
//}