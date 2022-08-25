package fileCenter

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"im-center/common/dynamicConfig/cfgByWatchFile"
	"im-center/common/tool"
	"strconv"
	"strings"
	"time"
)

// 查找重试的文件分片列表
func findRetrySeq(dirPath string, metadata *ServerFileMetadata) []int {
	slices := []int{}
	// 获取已保存的文件片序号
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("获取该文件夹里面文件列表出错:%s, err:%+v\n", dirPath, err)
		return nil
	}
	//storeSeq := make(map[string]bool)
	//去除有问题的分片，记录已经上传的分片
	//for _, file := range files {
	//	seq := file.Name()
	//	if file.IsDir() || seq[:1] == "." {
	//		fmt.Printf("文件片丢弃::%s\n", seq)
	//		continue
	//	}
	//	if metadata.FileDataType == "dicom" || strings.HasSuffix(file.Name(), ".dcm") {
	//		seq = strings.Split(file.Name(), ".dcm")[0]
	//		if _, err := strconv.Atoi(seq); err != nil {
	//			fmt.Printf("该分片id不符合规则:%s, err:%+v\n", seq, err)
	//			continue
	//		}
	//	} else if _, err := strconv.Atoi(file.Name()); err != nil {
	//		fmt.Printf("文件片有错:%v, name:%s\n", err, file.Name())
	//		continue
	//	}
	//	storeSeq[seq] = true
	//}
	//
	//i := 0
	//for ; i < metadata.SliceNum && len(storeSeq) > 0; i++ {
	//	indexStr := strconv.Itoa(i)
	//	if _, ok := storeSeq[indexStr]; ok {
	//		delete(storeSeq, indexStr)
	//	} else {
	//		slices = append(slices, i)
	//	}
	//}

	storeSeq := make(map[int]bool)
	for _, file := range files {
		var seq int
		if file.IsDir() || file.Name()[:1] == "." {
			fmt.Printf("文件片丢弃::%s\n", file.Name())
			continue
		}
		if metadata.FileDataType == "dicom" || strings.HasSuffix(file.Name(), ".dcm") {
			if seq, err = strconv.Atoi(strings.Split(file.Name(), ".dcm")[0] ); err != nil {
				fmt.Printf("该分片id不符合规则:%s, err:%+v\n", seq, err)
				continue
			}
		} else if seq, err = strconv.Atoi(file.Name()); err != nil {
			fmt.Printf("文件片有错:%v, name:%s\n", err, file.Name())
			continue
		}
		storeSeq[seq] = true
	}
	if len(storeSeq) == 0 {
		for v := 0; v < metadata.SliceNum; v++ { slices = append(slices, v) }
	} else {
		for v := 0; v < metadata.SliceNum; v++ {
			if _, ok := storeSeq[v]; !ok {
				slices = append(slices, v)
			}
		}
	}

	// -1指代slices的最大数字序号到最后一片都没有收到
	//if i < metadata.SliceNum {
	//	slices = append(slices, i)
	//	i += 1
	//	if i < metadata.SliceNum {
	//		slices = append(slices, -1)
	//	}
	//}

	fmt.Println("还需重传的片", slices)
	return slices
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给文件是否存在
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil || s.IsDir() {
		//fmt.Println("文件不存在或为文件夹：", path)
		return false
	}
	return true
}

func CopyFile(from string, dst string) error {
	fromS, err := os.Stat(from)
	if err != nil {
		return err
	}
	if fromS.IsDir() {
		if _, err = os.Stat(dst); err != nil {
			if err = os.MkdirAll(dst, 0766); err != nil{
				return err
			}
		}
		list, err := ioutil.ReadDir(from)
		if err != nil {err = errors.New("原路径为空") }
		for _, item := range list{
			if err = CopyFile(path.Join(from, item.Name()), path.Join(dst, item.Name())); err != nil {
				//出现错误，把已经拷贝的文件删除
				os.RemoveAll(dst)
				return err
			}
		}
	} else {
		//读取原文件
		file, err := os.Open(from)
		if err != nil {
			return err
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		//创建目标文件
		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, bufReader)
		if err != nil {
			//出现错误，把已经拷贝的文件删除
			os.Remove(dst)
		}
	}
	return err
}

// 获取文件 Fid
func GetFid(nodeId string, md5sum string) string {
	t := time.Now().In(tool.Loc)
	return fmt.Sprintf("%s-%d%02d%02d-%s", nodeId, t.Year(), t.Month(), t.Day(), md5sum)
}

// 加载元数据文件信息
func LoadMetadata(filePath string) (*ServerFileMetadata, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("获取文件状态失败，文件路径：", filePath)
		return nil, err
	}
	var metadata ServerFileMetadata
	filedata := gob.NewDecoder(file)
	err = filedata.Decode(&metadata)
	if err != nil {
		fmt.Println("格式化文件元数据失败, err", err)
		return nil, err
	}
	return &metadata, nil
}

// 写元数据文件信息
func StoreMetadata(filePath string, metadata *ServerFileMetadata) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		fmt.Println("创建元数据文件失败")
		return err
	}

	enc := gob.NewEncoder(f)
	err = enc.Encode(metadata)
	if err != nil {
		fmt.Println("写元数据文件失败")
		return err
	}
	return nil
}

func CheckFileExist(node *cfgByWatchFile.Node, filePath string) {
	defer wg.Done()
	if ok := IsFile(fmt.Sprintf("%s/%s", node.DiskMountedOn, filePath)); ok {return }
	//s, err := os.Stat(fmt.Sprintf("%s/%s", node.DiskMountedOn, filePath))
	//if err != nil || s.IsDir() {
	//	fmt.Println("获取文件失败/文件为文件夹，文件路径：", filePath)
	//	return
	//}
	val, ok := m.LoadOrStore(fmt.Sprintf(existFileNodes, filePath), []*cfgByWatchFile.Node{node})
	vals := val.([]*cfgByWatchFile.Node)
	if ok {
		// 里面存在值
		vals = append(vals, node)
	}
	// 存储了值
}

func writeFile(src io.Reader, dst string) error {
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("写入文件失败:%+v", err)
		return err
	}
	defer f.Close()
	io.Copy(f, src)
	return nil
}

//通过SourceFileId获取Fid
func GetFidBySourceFileId(sourceFileId string) string {
	return strings.Split(sourceFileId, ";")[0]
}

func SliceIdCheckForDicom(slicesId string, metadata *ServerFileMetadata) (seq string, err error) {
	err = errors.New("不是dicom文件")
	if metadata.FileDataType == "dicom" && strings.HasSuffix(slicesId, ".dcm") {
		seq = strings.Split(slicesId, ".dcm")[0]
		if _, err = strconv.Atoi(seq); err != nil {
			seq = ""
			fmt.Printf("该分片id不符合规则:%s, err:%+v\n", slicesId, err)
		}
	}
	return seq, err
}

//// GetMetadataFilepath 获取文件元数据文件路径
//func GetMetadataFilepath(filePath string) string {
//	return filePath + ".slice"
//}
//// 通过filename加载元数据并比对fid，确认文件是否存在
//func CheckFileExistByMetadata(fid string, filename string, storeDir string) (bool, error) {
//	metadataPath := GetMetadataFilepath(path.Join(storeDir, filename))
//	// 加载元数据信息
//	metadata, err := LoadMetadata(metadataPath)
//	if err != nil {
//		return false, err
//	}
//	// 校验fid和filename是匹配的
//	if metadata.Fid != fid {
//		fmt.Println("文件名和fid对不上，请确认后重试")
//		return false, errors.New("文件名和fid对不上，请确认后重试")
//	}
//	return true, nil
//}

func mergeSlice(dstPath string, srcPath string, sliceNum int) error {
	if IsFile(dstPath) {
		fmt.Print("目标文件已存在")
		return nil
	}
	zf, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer zf.Close()
	for i := 0; i < sliceNum; i++ {
		sliceFile, err := os.Open(path.Join(srcPath, strconv.Itoa(i)))
		if err != nil {
			return err
		}
		// 偏移量需要重新进行调整
		sliceFile.Seek(0, 0)
		defer sliceFile.Close()
		_, err = io.Copy(zf, sliceFile)
		if err != nil {
			return err
		}
	}
	return nil
}